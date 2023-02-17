package binance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
	"net/http"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
	"strings"
)

type Platform struct {
	platform.PlatformTemplate
	Binance *binance.Client
}

func New(name string, url string, tradeTypes []string, payTypes []string, tokens []string, allTokens []string) *Platform {
	return &Platform{
		PlatformTemplate: platform.PlatformTemplate{
			Name:       name,
			Url:        url,
			TradeTypes: tradeTypes,
			PayTypes:   payTypes,
			Tokens:     tokens,
			AllTokens:  allTokens,
			Client:     http.Client{}},
		Binance: binance.NewClient("", ""),
	}

}

func (p *Platform) GetResult(c *config.Configuration) (*platform.ResultPlatformData, error) {
	result := platform.ResultPlatformData{}
	spotData, err := p.GetSpotData()
	if err != nil {
		return nil, fmt.Errorf("can't get spot data: %w", err)
	}

	result.Name = p.Name
	result.Spot = *spotData

	for _, token := range p.Tokens {
		buy, err := p.GetAdvertise(c, token, p.TradeTypes[0])
		sell, err := p.GetAdvertise(c, token, p.TradeTypes[1])
		if err != nil {
			return nil, fmt.Errorf("can't get advertise: %w", err)
		}
		result.Tokens[token] = platform.TokenInfo{
			Buy:  *buy,
			Sell: *sell,
		}

	}
	return &result, nil
}

func (p *Platform) GetAdvertise(c *config.Configuration, token string, tradeType string) (*platform.Advertise, error) {
	userConfig := &c.UserConfig

	query, err := p.GetQuery(userConfig, token, tradeType)
	if err != nil {
		return nil, fmt.Errorf("can't get query: %w", err)
	}

	response, err := p.DoRequest(query)
	if err != nil {
		return nil, fmt.Errorf("can't do request to get binance response: %w", err)
	}

	advertise, err := p.ResponseToAdvertise(&response)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}

	log.Printf("advertise: %+v", advertise)

	return advertise, nil
}

func (p *Platform) GetSpotData() (*map[string]float64, error) {
	allTokens := p.AllTokens
	set := *p.CreatePairsSet(allTokens)
	result := map[string]float64{}

	prices, err := p.Binance.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't get spot data: %w", err)
	}
	for _, p := range prices {
		if set[p.Symbol] {
			result[p.Symbol], _ = strconv.ParseFloat(p.Price, 64)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("can't parse price to float: %w", err)
	}
	return &result, nil
}

func (p *Platform) GetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {

	var BinanceJsonData = map[string]interface{}{
		"proMerchantAds": false,
		"page":           1,
		"rows":           10,
		"payTypes":       p.GetPayTypes(c),
		"countries":      []string{},
		"publisherType":  nil,
		"transAmount":    c.MinValue,
		"asset":          token,
		"fiat":           "RUB",
		"tradeType":      tradeType, /*BUY(Купить) or SELL(продать)*/
	}

	result, err := p.QueryToBytes(&BinanceJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

func (p *Platform) ResponseToAdvertise(response *[]byte) (*platform.Advertise, error) {
	var data Response

	err := json.Unmarshal(*response, &data)
	if err != nil {
		return nil, fmt.Errorf("cant' unmarshall data from binance response: %w", err)
	}
	item := data.Data[0]
	cost, _ := strconv.ParseFloat(item.Adv.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.Adv.MinSingleTransAmount, 64)
	maxLimit, _ := strconv.ParseFloat(item.Adv.MaxSingleTransAmount, 64)
	available, _ := strconv.ParseFloat(item.Adv.DynamicMaxSingleTransQuantity, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.Advertiser.NickName,
		Asset:        item.Adv.Asset,
		Fiat:         item.Adv.FiatUnit,
		BankName:     payTypesToString(&data),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.Advertiser.MonthOrderCount,
		TradeType:    item.Adv.TradeType,
		Available:    available,
	}, nil
}

func payTypesToString(r *Response) string {
	data := r.Data[0].Adv.TradeMethods
	var result []string
	for _, k := range data {
		result = append(result, k.TradeMethodName)

	}
	return strings.Join(result, ", ")
}
