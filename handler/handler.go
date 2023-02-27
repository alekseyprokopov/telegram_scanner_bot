package handler

import (
	"log"
	"scanner_bot/config"
	"scanner_bot/platform"
	"scanner_bot/platform/binance"
	"scanner_bot/platform/bybit"
	"scanner_bot/platform/huobi"
	"sync"
	"time"
)

type PlaftormHandler struct {
	Platforms map[string]platform.Platform
	//Binance platform.Platform
	//Bybit   platform.Platform
	//Huobi   platform.Platform
}

func New() *PlaftormHandler {
	return &PlaftormHandler{
		Platforms: map[string]platform.Platform{
			Binance.name: binance.New(Binance.name, Binance.p2pURL, Binance.tradeTypes, Binance.tokens, Binance.tokensDict, Binance.payTypesDict, allPairs),
			Bybit.name:   bybit.New(Bybit.name, Bybit.p2pURL, Bybit.tradeTypes, Bybit.tokens, Bybit.tokensDict, Bybit.payTypesDict, allPairs),
			Huobi.name:   huobi.New(Huobi.name, Huobi.p2pURL, Huobi.tradeTypes, Huobi.tokens, Huobi.tokensDict, Huobi.payTypesDict, allPairs),
		}}
	//Binance: binance.New(Binance.name, Binance.p2pURL, Binance.tradeTypes, Binance.tokens, Binance.tokensDict, Binance.payTypesDict, allPairs),
	//Bybit:   bybit.New(Bybit.name, Bybit.p2pURL, Bybit.tradeTypes, Bybit.tokens, Bybit.tokensDict, Bybit.payTypesDict, allPairs),
	//Huobi:   huobi.New(Huobi.name, Huobi.p2pURL, Huobi.tradeTypes, Huobi.tokens, Huobi.tokensDict, Huobi.payTypesDict, allPairs)}

}

func (p *PlaftormHandler) InsideTakerTaker(c *config.Configuration) *[]Chain {
	start := time.Now()
	var chains []Chain //data, err := p.Huobi.GetResult(c)
	wg := sync.WaitGroup{}
	for key, value := range p.Platforms {
		key, value := key, value
		wg.Add(1)
		go func() {
			platformResult, err := value.GetResult(c)
			if err != nil {
				log.Printf("\ncan't get result from: %s\n", key)
			}
			p.findInsideTT(platformResult, &chains)
			defer wg.Done()
		}()

	}
	log.Println("TIME : ", time.Since(start))
	wg.Wait()
	return &chains
}

func IsExistPair(pair string, data *platform.ResultPlatformData) bool {
	_, ok := data.Spot[pair]
	return ok
}

func (p *PlaftormHandler) findInsideTT(data *platform.ResultPlatformData, chains *[]Chain) {
	for token1, tokenInfo1 := range data.Tokens {
		for token2, tokenInfo2 := range data.Tokens {
			if token1 == token2 {
				continue
			}
			pair1name := token1 + token2
			pair1spotPrice, ok1 := data.Spot[pair1name]

			pair2name := token2 + token1
			pair2spotPrice, ok2 := data.Spot[pair2name]

			var result float64
			var spotPrice float64
			var pairName string

			if ok1 && tokenInfo1.Buy.Cost != 0 && tokenInfo2.Sell.Cost != 0 {
				result = 100/tokenInfo1.Buy.Cost*pair1spotPrice*tokenInfo2.Sell.Cost - 100
				spotPrice = pair1spotPrice
				pairName = pair1name
			}

			if ok2 && tokenInfo1.Buy.Cost != 0 && tokenInfo2.Sell.Cost != 0 {
				result = 100/tokenInfo1.Buy.Cost/pair2spotPrice*tokenInfo2.Sell.Cost - 100
				spotPrice = pair2spotPrice
				pairName = pair2name

			}

			if result > 0 {
				chain := Chain{
					PairName:  pairName,
					Buy:       &tokenInfo1.Buy,
					Sell:      &tokenInfo2.Sell,
					SpotPrice: spotPrice,
					Profit:    result,
				}
				*chains = append(*chains, chain)
			}
		}
	}

}

type Chain struct {
	Buy       *platform.Advertise
	Sell      *platform.Advertise
	PairName  string
	SpotPrice float64
	Profit    float64
}
