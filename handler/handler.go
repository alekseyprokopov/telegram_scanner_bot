package handler

import (
	"log"
	"scanner_bot/config"
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

func (p *PlaftormHandler) InsideTakerTaker(c *config.Configuration) {
	data, err := p.Binance.GetResult(c)
	if err != nil {
		log.Println("err")
	}
	for token1, tokenInfo1 := range data.Tokens {
		for token2, tokenInfo2 := range data.Tokens {
			if token1 == token2 {
				continue
			}
			pair1, ok1 := data.Spot[token1+token2]
			pair2, ok2 := data.Spot[token1+token2]

			if ok1 {
				result := 1 / tokenInfo1.Buy.Cost * pair1 / tokenInfo2.Sell.Cost
				log.Println("\nПАРА: ", token1+token2)
				log.Printf("ПОКУПКА %s: %f \n", token1, tokenInfo1.Buy.Cost)
				log.Printf("СПОТ: %f \n", pair1)
				log.Printf("ПОДАЖА %s: %f\n", token2, tokenInfo2.Sell.Cost)
				log.Println("ПРОФИТ: ", result)
				log.Println("-------------")

			}

			if ok2 {
				result := 1 / tokenInfo1.Buy.Cost * pair1 / tokenInfo2.Sell.Cost
				log.Println("ПАРА: ", token2+token1)
				log.Printf("\nПОКУПКА %s: %f \n", token1, tokenInfo1.Buy.Cost)
				log.Println("СПОТ: ", pair2)
				log.Printf("\nПОДАЖА %s: %f\n", tokenInfo2.Sell.Cost)
				log.Println("ПРОФИТ: ", result)
			}
		}
	}

}

func IsExistPair(pair string, data *platform.ResultPlatformData) bool {
	_, ok := data.Spot[pair]
	return ok
}
