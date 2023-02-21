package handler

import (
	"scanner_bot/platform"
	"scanner_bot/platform/bybit"
)

type PlaftormHandler struct {
	Binance platform.Platform
	Bybit   platform.Platform
	Huobi   platform.Platform
}

var (
	//name             = "binance"
	//url              = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	//payTypes         = []string{"RosBankNew", "TinkoffNew"}
	//tokens           = []string{"USDT", "BTC", "BUSD", "BNB", "ETH", "SHIB"}
	//BinanceTradeType = []string{"BUY", "SELL"}

	name     = "bybit"
	url      = "https://api2.bybit.com/fiat/otc/item/online"
	payTypes = []string{"185", "75", "64", "62", "274"}
	tokens   = []string{"USDT", "BTC", "ETH", "USDC"}

	BybitTradeType = []string{"1", "0"} //side 1- купить. 0 - продать
)

func New() *PlaftormHandler {
	return &PlaftormHandler{Bybit: bybit.New(name, url, BybitTradeType, payTypes, tokens, tokens)}
}
