package bybit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
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
	return p.TemplateResult(c, p.spotData, p.advertise)
}

func (p *Platform) spotData() (*map[string]float64, error) {
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

	for _, item := range spotResponse.Result.List {
		_, ok := set[item.Symbol]
		if ok {
			price, err := strconv.ParseFloat(item.LastPrice, 64)
			if err != nil {
				return nil, fmt.Errorf("can't parse to Float: %w", err)
			}

			result[item.Symbol] = price
		}
	}
	return &result, err

}

func (p *Platform) advertise(c *config.Configuration, token string, tradeType string) (*platform.Advertise, error) {
	userConfig := &c.UserConfig

	query, err := p.getQuery(userConfig, token, tradeType)
	if err != nil {
		return nil, fmt.Errorf("can't get query: %w", err)
	}
	response, err := p.DoPostRequest(query)
	if err != nil {
		return nil, fmt.Errorf("can't do request to get bybit response: %w", err)
	}

	advertise, err := p.responseToAdvertise(response, userConfig)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}
	return advertise, nil
}

func (p *Platform) getQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var BybitJsonData = map[string]interface{}{
		"userId":     "",
		"tokenId":    token,
		"currencyId": "RUB",
		"payment":    p.GetPayTypes(c),
		"side":       tradeType,
		"size":       "10",
		"page":       "1",
		"amount":     strconv.Itoa(c.MinValue),
	}
	result, err := p.QueryToBytes(&BybitJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

func (p *Platform) responseToAdvertise(response *[]byte, config *config.Config) (*platform.Advertise, error) {
	orders := config.Orders
	var data Response
	err := json.Unmarshal(*response, &data)

	if err != nil || len(data.Result.Items) == 0 {
		return nil, fmt.Errorf("cant' unmarshall data from  response: %w", err)
	}
	var adv *AdvertiseData

	for _, item := range data.Result.Items {
		if item.RecentOrderNum >= orders {
			adv = &item
			break
		}
	}
	if adv == nil {
		adv = &data.Result.Items[0]
	}

	cost, _ := strconv.ParseFloat(adv.Price, 64)
	minLimit, _ := strconv.ParseFloat(adv.MinAmount, 64)
	maxLimit, _ := strconv.ParseFloat(adv.MaxAmount, 64)
	available, _ := strconv.ParseFloat(adv.LastQuantity, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   adv.NickName,
		Asset:        adv.TokenID,
		Fiat:         adv.CurrencyID,
		BankName:     p.PayTypesToString(adv.Payments),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  adv.RecentOrderNum,
		TradeType:    bybitTradeType(adv.Side),
		Available:    available,
	}, nil
}

func bybitTradeType(i int) string {
	if i == 1 {
		return "BUY"
	}
	return "SELL"
}
