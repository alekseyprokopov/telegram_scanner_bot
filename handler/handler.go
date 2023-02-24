package handler

import (
	"scanner_bot/platform"
	"scanner_bot/platform/binance"
	"scanner_bot/platform/bybit"
	"scanner_bot/platform/huobi"
)

type PlaftormHandler struct {
	Binance platform.Platform
	Bybit   platform.Platform
	Huobi   platform.Platform
}

func New() *PlaftormHandler {
	return &PlaftormHandler{
		Binance: binance.New(Binance.name, Binance.p2pURL, Binance.tradeTypes, Binance.tokens, Binance.tokensDict, Binance.payTypesDict, allPairs),
		Bybit:   bybit.New(Bybit.name, Bybit.p2pURL, Bybit.tradeTypes, Bybit.tokens, Bybit.tokensDict, Bybit.payTypesDict, allPairs),
		Huobi:   huobi.New(Huobi.name, Huobi.p2pURL, Huobi.tradeTypes, Huobi.tokens, Huobi.tokensDict, Huobi.payTypesDict, allPairs)}

}
