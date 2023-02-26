package handler

import (
	"log"
	"scanner_bot/config"
	"scanner_bot/platform"
	"scanner_bot/platform/binance"
	"time"
)

type PlaftormHandler struct {
	Binance platform.Platform
	Bybit   platform.Platform
	Huobi   platform.Platform
}

func New() *PlaftormHandler {
	return &PlaftormHandler{
		Binance: binance.New(Binance.name, Binance.p2pURL, Binance.tradeTypes, Binance.tokens, Binance.tokensDict, Binance.payTypesDict, allPairs),
		//Bybit:   bybit.New(Bybit.name, Bybit.p2pURL, Bybit.tradeTypes, Bybit.tokens, Bybit.tokensDict, Bybit.payTypesDict, allPairs),
		//Huobi:   huobi.New(Huobi.name, Huobi.p2pURL, Huobi.tradeTypes, Huobi.tokens, Huobi.tokensDict, Huobi.payTypesDict, allPairs)}
	}
}

func (p *PlaftormHandler) InsideTakerTaker(c *config.Configuration) {
	start:= time.Now()
	data, err := p.Binance.GetResult(c)
	//log.Printf("DATA: %+v ", data)
	//log.Printf("BUY: %+v ", data.Tokens["USDT"].Buy)
	log.Println("TIME : ", time.Since(start))

	if err != nil {
		log.Println("err")
	}
	for token1, tokenInfo1 := range data.Tokens {
		for token2, tokenInfo2 := range data.Tokens {
			if token1 == token2 {
				continue
			}
			pair1, ok1 := data.Spot[token1+token2]
			pair2, ok2 := data.Spot[token2+token1]

			if ok1 {
				log.Println("-------------")
				log.Println("pair1")
				result := 1/tokenInfo1.Buy.Cost * pair1 * tokenInfo2.Sell.Cost
				log.Println("ПАРА: ", token1+token2)
				log.Printf("ПОКУПКА %s: %f \n", token1, tokenInfo1.Buy.Cost)
				log.Printf("СПОТ: %f \n", pair1)
				log.Printf("ПОДАЖА %s: %f\n", token2, tokenInfo2.Sell.Cost)
				log.Println("ПРОФИТ: ", result)
				log.Println("-------------")

			}

			if ok2 {
				log.Println("-------------")

				log.Println("pair2")
				result := 1/tokenInfo1.Buy.Cost / pair2 * tokenInfo2.Sell.Cost
				log.Println("ПАРА: ", token2+token1)
				log.Printf("ПОКУПКА %s: %f \n", token1, tokenInfo1.Buy.Cost)
				log.Println("СПОТ: ", pair2)
				log.Printf("ПОДАЖА %s: %f\n", token2, tokenInfo2.Sell.Cost)
				log.Println("ПРОФИТ: ", result)
			}
		}
	}

}

func IsExistPair(pair string, data *platform.ResultPlatformData) bool {
	_, ok := data.Spot[pair]
	return ok
}
