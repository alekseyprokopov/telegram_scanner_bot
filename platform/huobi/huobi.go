package huobi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
	"strings"
	"sync"
)

type Platform struct {
	*platform.PlatformTemplate
}

func New(name string, url string, apiUrl string, tradeTypes []string, tokens []string, tokensDict map[string]string, payTypesDict map[string]string, allPairs map[string]bool) *Platform {
	return &Platform{
		PlatformTemplate: platform.New(name, url, apiUrl, tradeTypes, tokens, tokensDict, payTypesDict, allPairs),
	}
}

func (p *Platform) GetResult(c *config.Configuration) (*platform.ResultPlatformData, error) {
	result := platform.ResultPlatformData{}
	result.Tokens = map[string]*platform.TokenInfo{}
	wg := sync.WaitGroup{}
	var mu sync.Mutex
	result.Name = p.Name

	wg.Add(1)
	go func() {
		spotData, err := p.getSpotData()
		if err != nil {
			log.Printf("can't get spot data: %v", err)
		}
		mu.Lock()
		result.Spot = *spotData
		mu.Unlock()
		defer wg.Done()
	}()

	for _, token := range p.Tokens {
		token := token
		//for huobi tokens
		realToken := p.TokensDict[token]
		//
		tokenInfo := &platform.TokenInfo{}
		result.Tokens[realToken] = tokenInfo

		wg.Add(1)
		go func() {
			buy, err := p.getAdvertise(c, token, p.TradeTypes[0])
			if err != nil || buy == nil {
				log.Printf("can't get buy advertise for huobi, token (%s): %v", token, err)
			} else {
				mu.Lock()
				tokenInfo.Buy = *buy
				mu.Unlock()
			}
			defer wg.Done()
		}()

		wg.Add(1)
		go func() {
			sell, err := p.getAdvertise(c, token, p.TradeTypes[1])
			if err != nil || sell == nil {
				log.Printf("can't get sell advertise for huobi, token (%s): %v", token, err)
			} else {
				mu.Lock()
				tokenInfo.Sell = *sell
				mu.Unlock()
			}
			defer wg.Done()
		}()
		//result.Tokens[token] = tokenInfo

	}
	wg.Wait()

	return &result, nil
}
func (p *Platform) getSpotData() (*map[string]float64, error) {
	data, err := p.DoGetRequest(p.ApiUrl, "")
	if err != nil {
		return nil, fmt.Errorf("can't do getRequest to huobi API: %w", err)
	}

	var spotResponse SpotResponse
	if err := json.Unmarshal(*data, &spotResponse); err != nil {
		return nil, fmt.Errorf("can't unmarshall: %w", err)
	}

	result := map[string]float64{}
	set := p.AllPairs

	for _, item := range spotResponse.Data {
		_, ok := set[strings.ToUpper(item.Symbol)]
		if ok {
			result[strings.ToUpper(item.Symbol)] = item.Close
		}
	}
	return &result, err

}

func (p *Platform) getAdvertise(c *config.Configuration, token string, tradeType string) (*platform.Advertise, error) {
	userConfig := &c.UserConfig
	query := p.getQuery(userConfig, token, tradeType)
	response, err := p.DoGetRequest(p.Url, query)
	if err != nil {
		return nil, fmt.Errorf("can't do request to get bybit response: %w", err)
	}
	advertise, err := p.responseToAdvertise(response, userConfig)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}
	return advertise, nil
}

func (p *Platform) getQuery(c *config.Config, token string, tradeType string) string {
	u := url.Values{
		"coinId":       []string{token},     //usdt
		"currency":     []string{"11"},      //rub
		"tradeType":    []string{tradeType}, //buy
		"currPage":     []string{"1"},
		"payMethod":    []string{strings.Join(p.GetPayTypes(c), ",")},
		"acceptOrder":  []string{"0"},
		"country":      []string{""},
		"blockType":    []string{"general"},
		"online":       []string{"1"},
		"range":        []string{"0"},
		"amount":       []string{c.MinValue}, //amount
		"isThumbsUp":   []string{"false"},
		"isMerchant":   []string{"false"},
		"isTraded":     []string{"false"},
		"onlyTradable": []string{"false"},
		"isFollowed":   []string{"false"},
	}

	return u.Encode()
}

func (p *Platform) responseToAdvertise(response *[]byte, config *config.Config) (*platform.Advertise, error) {
	orders := config.Orders
	var data Response
	err := json.Unmarshal(*response, &data)

	if err != nil || len(data.Data) == 0 || data.Code != 200 {
		return nil, fmt.Errorf("cant' unmarshall data from huobi response: %w", err)
	}
	var adv *AdvertiseData
	for _, item := range data.Data {
		if item.TradeMonthTimes >= orders {
			adv = &item
			break
		}
	}
	if adv == nil {
		adv = &data.Data[0]
	}


	cost, _ := strconv.ParseFloat(adv.Price, 64)
	minLimit, _ := strconv.ParseFloat(adv.MinTradeLimit, 64)
	maxLimit, _ := strconv.ParseFloat(adv.MaxTradeLimit, 64)
	available, _ := strconv.ParseFloat(adv.TradeCount, 64)
	pays := getStringSlice(adv.PayMethods)

	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   adv.UserName,
		Asset:        p.TokenFromDict(strconv.Itoa(adv.CoinID)),
		Fiat:         fiatDict[adv.Currency],
		BankName:     p.PayTypesToString(pays),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  adv.TradeMonthTimes,
		TradeType:    huobiTradeType(adv.TradeType),
		Available:    available,
	}, nil
}

var fiatDict = map[int]string{
	11: "RUB",
}
