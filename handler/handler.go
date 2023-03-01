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
func (p *PlaftormHandler) OutsideTakerTaker(c *config.Configuration) *[]Chain {
	start := time.Now()
	var chains []Chain //data, err := p.Huobi.GetResult(c)
	var platformResults []*platform.ResultPlatformData
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	for key, value := range p.Platforms {
		key, value := key, value
		wg.Add(1)
		go func() {
			platformResult, err := value.GetResult(c)
			if err != nil {
				log.Printf("\ncan't get result from: %s\n", key)
			}
			mu.Lock()
			platformResults = append(platformResults, platformResult)
			mu.Unlock()
			defer wg.Done()
		}()
	}
	log.Println("TIME : ", time.Since(start))
	wg.Wait()
	log.Printf("%s, %+v", platformResults[0].Name, platformResults[0].Spot)
	log.Printf("%s, %+v", platformResults[1].Name, platformResults[1].Spot)
	log.Printf("%s, %+v", platformResults[2].Name, platformResults[2].Spot)

	for _, item1 := range platformResults {
		for _, item2 := range platformResults {
			p.findOutsideTTspot1(item1, item2, &chains)
			p.findOutsideTTspot2(item1, item2, &chains)
		}
	}
	//for i := 0; i < len(platformResults)-1; i++ {
	//	for j := i + 1; j < len(platformResults); j++ {
	//		p.findOutsideTT(platformResults[i], platformResults[j], &chains)
	//	}
	//}

	log.Println("TIME : ", time.Since(start))
	return &chains
}
func (p *PlaftormHandler) findOutsideTTspot1(first *platform.ResultPlatformData, second *platform.ResultPlatformData, chains *[]Chain) {
	//
	counter := 0
	for token1, tokenInfo1 := range first.Tokens {
		for token2, tokenInfo2 := range second.Tokens {
			if token1 == token2 {
				continue
			}
			pair1name := token1 + token2
			pair1spotPrice, ok1 := first.Spot[pair1name]

			pair2name := token2 + token1
			pair2spotPrice, ok2 := first.Spot[pair2name]

			var result float64
			var spotPrice float64
			var pairName string
			//
			counter += 1
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
			//log.Println("----------------")
			//
			//log.Println("RESULT: ", result)
			//log.Println("ok1: ", ok1)
			//log.Println("ok1: ", ok2)
			//
			//log.Println("----------------")

			if result > 0 {
				chain := &Chain{
					PairName:  pairName,
					Buy:       &tokenInfo1.Buy,
					Sell:      &tokenInfo2.Sell,
					SpotName:  tokenInfo1.Buy.PlatformName,
					SpotPrice: spotPrice,
					Profit:    result,
				}
				*chains = append(*chains, *chain)
			}
		}
	}
	log.Println("COUNTER: ", counter, " NAME: ", first.Name, " Name: ", second.Name)
}

func (p *PlaftormHandler) findOutsideTTspot2(first *platform.ResultPlatformData, second *platform.ResultPlatformData, chains *[]Chain) {
	//
	counter := 0
	for token1, tokenInfo1 := range first.Tokens {
		for token2, tokenInfo2 := range second.Tokens {
			if token1 == token2 {
				continue
			}
			pair1name := token1 + token2
			pair1spotPrice, ok1 := second.Spot[pair1name]

			pair2name := token2 + token1
			pair2spotPrice, ok2 := second.Spot[pair2name]

			var result float64
			var spotPrice float64
			var pairName string
			//
			counter += 1
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

			if result != 0 {
				chain := &Chain{
					PairName:  pairName,
					Buy:       &tokenInfo1.Buy,
					Sell:      &tokenInfo2.Sell,
					SpotName:  tokenInfo2.Sell.PlatformName,
					SpotPrice: spotPrice,
					Profit:    result,
				}
				*chains = append(*chains, *chain)
			}
		}
	}
	log.Println("COUNTER: ", counter, " NAME: ", first.Name, " Name: ", second.Name)
}

func (p *PlaftormHandler) InsideTakerTaker(c *config.Configuration) *[]Chain {
	start := time.Now()
	var chains []Chain //data, err := p.Huobi.GetResult(c)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for key, value := range p.Platforms {
		key, value := key, value
		wg.Add(1)
		go func() {
			platformResult, err := value.GetResult(c)
			//log.Printf("platformRESULT SPOT:%s : %+v ", platformResult.Name, platformResult.Spot)
			if err != nil {
				log.Printf("\ncan't get result from: %s\n", key)
			}
			mu.Lock()
			p.findInsideTT(platformResult, &chains)
			mu.Unlock()
			defer wg.Done()
		}()
	}
	wg.Wait()
	log.Println("TIME : ", time.Since(start))
	//log.Println(len(chains))
	return &chains
}

func IsExistPair(pair string, data *platform.ResultPlatformData) bool {
	_, ok := data.Spot[pair]
	return ok
}

func (p *PlaftormHandler) findInsideTT(data *platform.ResultPlatformData, chains *[]Chain) {
	log.Printf("PLATFORMRESULT %s, %+v\n", data.Name, data.Spot)
	log.Printf("PLATFORMTOKENS %s, %+v\n", data.Name, data.Tokens)
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

			//if result > -100 {
			chain := Chain{
				PairName:  pairName,
				Buy:       &tokenInfo1.Buy,
				Sell:      &tokenInfo2.Sell,
				SpotPrice: spotPrice,
				Profit:    result,
			}
			*chains = append(*chains, chain)
			//}
		}
	}

}

type Chain struct {
	Buy       *platform.Advertise
	Sell      *platform.Advertise
	PairName  string
	SpotPrice float64
	SpotName  string
	Profit    float64
}

//func (p *PlaftormHandler) OutsideTakerTaker(c *config.Configuration) *[]Chain {
//	start := time.Now()
//	var chains []Chain //data, err := p.Huobi.GetResult(c)
//	wg := sync.WaitGroup{}
//	checkCollision := map[string]bool{}
//	for key1, value1 := range p.Platforms {
//		key1, value1 := key1, value1
//
//		for key2, value2 := range p.Platforms {
//			//создание пар для рутин
//			key2, value2 := key2, value2
//			checkRepeat, _ := checkCollision[key2]
//			//исключение одинаковых бирж
//			if key1 == key2 && checkRepeat {
//				log.Print("HELLO ", key1, key2)
//				continue
//			}
//
//			wg.Add(1)
//			go func() {
//				platformResult1, err := value1.GetResult(c)
//				platformResult2, err := value2.GetResult(c)
//				if err != nil {
//					log.Printf("\ncan't get result from: %s, %s\n", key1, key2)
//				}
//				p.findOutsideTT(platformResult1, platformResult2, &chains)
//				defer wg.Done()
//			}()
//		}
//		checkCollision[key1] = true
//	}
//	log.Println("TIME : ", time.Since(start))
//	wg.Wait()
//	return &chains
//}
