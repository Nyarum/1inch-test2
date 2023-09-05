package poller

import (
	"1inch/cmd/service/uniswap"
	"1inch/internal"
	"1inch/internal/contracts"
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

type ExternalQueue interface {
	Push(path string, msg internal.Message) error
	Pull(path string) <-chan internal.Message
}

type Poller struct {
	ethClient            *ethclient.Client
	externalQueueClient  ExternalQueue
	msgPricesCh          chan internal.MsgPrices
	currentPricesOnPools map[string]internal.MsgPrices
}

func NewPoller(
	ethClient *ethclient.Client,
	externalQueueClient ExternalQueue,
	uniswapV2Pools []string,
	uniswapV3Pools []string,
) *Poller {
	ctx := context.Background()

	poller := &Poller{
		ethClient:            ethClient,
		externalQueueClient:  externalQueueClient,
		msgPricesCh:          make(chan internal.MsgPrices, 1),
		currentPricesOnPools: make(map[string]internal.MsgPrices),
	}

	go poller.watchUniswapV2(ctx, uniswapV2Pools)
	go poller.watchUniswapV3(ctx, uniswapV3Pools)

	return poller
}

func (p *Poller) Run() {
	for msg := range p.msgPricesCh {
		p.currentPricesOnPools[msg.Pool.Hex()] = msg

		log.Printf("Pushed a new update of prices to external queue: %v\n", msg)

		err := p.externalQueueClient.Push(internal.DefaultPath, internal.Message{
			Prices: msg,
		})
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (p *Poller) watchUniswapV3(ctx context.Context, addresses []string) {
	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.ContractsABI)))
	if err != nil {
		log.Println(err)
		return
	}

	query := ethereum.FilterQuery{}
	for _, addr := range addresses {
		query.Addresses = append(query.Addresses, common.HexToAddress(addr))
	}

	logs := make(chan types.Log)

	sub, err := p.ethClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-sub.Err():
			log.Println(err)
			return
		case vLog := <-logs:
			if vLog.Topics[0].Hex() != contractAbi.Events["Swap"].ID.Hex() {
				continue
			}

			var swapInfo struct {
				Amount0      *big.Int
				Amount1      *big.Int
				SqrtPriceX96 *big.Int
				Liquidity    *big.Int
				Tick         *big.Int
			}

			err := contractAbi.UnpackIntoInterface(&swapInfo, "Swap", vLog.Data)
			if err != nil {
				log.Println(err)
			}

			price0 := uniswap.SqrtPriceX96ToPrice(swapInfo.SqrtPriceX96, true)
			price1 := uniswap.SqrtPriceX96ToPrice(swapInfo.SqrtPriceX96, false)

			p.msgPricesCh <- internal.MsgPrices{
				Pool:   vLog.Address,
				Token0: price0,
				Token1: price1,
			}
		}
	}
}

func (p *Poller) watchUniswapV2(ctx context.Context, addresses []string) {
	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.Contractsv2pairABI)))
	if err != nil {
		log.Println(err)
		return
	}

	query := ethereum.FilterQuery{}
	for _, addr := range addresses {
		query.Addresses = append(query.Addresses, common.HexToAddress(addr))
	}

	logs := make(chan types.Log)

	sub, err := p.ethClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-sub.Err():
			log.Println(err)
			return
		case vLog := <-logs:
			if vLog.Topics[0].Hex() != contractAbi.Events["Swap"].ID.Hex() {
				continue
			}

			v2PairContract, err := contracts.NewContractsv2pair(vLog.Address, p.ethClient)
			if err != nil {
				log.Println(err)
				continue
			}

			reserves, err := v2PairContract.GetReserves(&bind.CallOpts{})
			if err != nil {
				log.Println(err)
				continue
			}

			p.msgPricesCh <- internal.MsgPrices{
				Pool:   vLog.Address,
				Token0: decimal.NewFromBigInt(reserves.Reserve0, 0),
				Token1: decimal.NewFromBigInt(reserves.Reserve1, 0),
			}
		}
	}
}
