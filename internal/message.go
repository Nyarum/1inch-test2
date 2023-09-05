package internal

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

const (
	DefaultPath = "msg"
)

type MsgPrices struct {
	Pool   common.Address
	Token0 decimal.Decimal
	Token1 decimal.Decimal
}

type Message struct {
	Prices MsgPrices
}
