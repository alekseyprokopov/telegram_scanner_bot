package handler

import (
	"scanner_bot/platform"
	"scanner_bot/platform/huobi"
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

	name     = "huobi"
	url      = "https://otc-akm.huobi.com/v1/data/trade-market"
	payTypes = []string{"185", "75", "64", "62", "274"}
	tokens   = []string{
		"2",  /*USDT*/
		"1",  /*BTC*/
		"62", /*USDD*/
		"4",  /*HT*/
		"22", /*TRX*/
		"3",  /*ETH*/
		"7",  /*XRP*/
		"8",  /*LTC*/
	}

	BybitTradeType = []string{"sell", "buy"} //side 1- купить. 0 - продать
)

func New() *PlaftormHandler {
	return &PlaftormHandler{Huobi: huobi.New(name, url, BybitTradeType, payTypes, tokens, tokens)}
}
