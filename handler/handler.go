package handler

import (
	"scanner_bot/platform"
	"scanner_bot/platform/binance"
)

type PlaftormHandler struct {
	Binance platform.Platform
}

var (
	name             = "binance"
	url              = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	payTypes         = []string{"RosBankNew", "TinkoffNew"}
	tokens           = []string{"USTD", "BTC", "BUSD", "BNB", "ETH", "SHIB"}
	BinanceTradeType = []string{"BUY", "SELL"}
)

func New() *PlaftormHandler {
	return &PlaftormHandler{Binance: binance.New(name, url, BinanceTradeType, payTypes, tokens, tokens)}
}
