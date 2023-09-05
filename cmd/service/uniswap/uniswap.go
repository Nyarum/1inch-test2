package uniswap

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
)

var (
	X96 = decimal.NewFromInt(2).Pow(decimal.NewFromInt(96))
)

// SqrtPriceX96ToPrice convert uniswap v3 sqrt price in x96 format to decimal.Decimal
// zeroForOne true: price = token0/token1 false: price = token1/token0
func SqrtPriceX96ToPrice(sqrtPriceX96 *big.Int, zeroForOne bool) (price decimal.Decimal) {
	d := decimal.NewFromBigInt(sqrtPriceX96, 0).Div(X96)
	p := d.Mul(d)

	if !zeroForOne {
		price = decimal.NewFromInt(1).Div(p)
		return
	}
	price = p
	return
}

func WeiToDecimal(wei *big.Float, delim *big.Float) *big.Float {
	return new(big.Float).Quo(wei, delim)
}

func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}
