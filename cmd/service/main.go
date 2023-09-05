package main

import (
	"1inch/cmd/service/poller"
	"1inch/internal"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nats-io/nats.go"
)

const (
	UniswapV2PoolAddrs = "0x2cc846fff0b08fb3bffad71f53a60b4b6e6d6482,0x81a460ea6fd96a73d5672f1f4aa684697d4b44cc,0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852"
	UniswapV3PoolAddrs = "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	uniswapv3Pools := flag.String("uniswapv3-pools", UniswapV3PoolAddrs, "Uniswap V3 addresses of pools")
	uniswapv2Pools := flag.String("uniswapv2-pools", UniswapV2PoolAddrs, "Uniswap V2 addresses of pools")
	flag.Parse()

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/06eaf9e210cd4587a85c1666dd1b2c17")
	if err != nil {
		log.Fatal(err)
	}

	natsClient, err := internal.NewNATS(natsURL)
	if err != nil {
		log.Fatal(err)
	}

	uniswapv2PoolsArr := strings.Split(*uniswapv2Pools, ",")
	uniswapv3PoolsArr := strings.Split(*uniswapv3Pools, ",")

	poll := poller.NewPoller(client, natsClient, uniswapv2PoolsArr, uniswapv3PoolsArr)
	poll.Run()
}
