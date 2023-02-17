package binance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
	"strings"
)

type Platform struct {
	platform.PlatformType
}

func New(name string, url string, payTypes []string, tokens []string, allTokens []string) *Platform {
	return &Platform{
		platform.PlatformType{
			Name:           name,
			Url:            url,
			PayTypes:       payTypes,
			PlatformTokens: tokens,
			AllTokents:     allTokens,
			Client:         http.Client{}},
	}
}

func (p *Platform) GetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {

	var BinanceJsonData = map[string]interface{}{
		"proMerchantAds": false,
		"page":           1,
		"rows":           10,
		"payTypes":       platform.GetPayTypes(platform.BinanceName, c),
		"countries":      []string{},
		"publisherType":  nil,
		"transAmount":    c.MinValue,
		"asset":          token,
		"fiat":           "RUB",
		"tradeType":      tradeType, /*BUY(Купить) or SELL(продать)*/
	}

	result, err := platform.QueryToBytes(&BinanceJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

func ResponseToAdvertise(response *[]byte) (*platform.Advertise, error) {
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
		PlatformName: platform.BinanceName,
		SellerName:   item.Advertiser.NickName,
		Asset:        item.Adv.Asset,
		Fiat:         item.Adv.FiatUnit,
		BankName:     binancePayTypesToString(&data),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.Advertiser.MonthOrderCount,
		TradeType:    item.Adv.TradeType,
		Available:    available,
	}, nil
}

func binancePayTypesToString(r *Response) string {
	data := r.Data[0].Adv.TradeMethods
	var result []string
	for _, k := range data {
		result = append(result, k.TradeMethodName)

	}
	return strings.Join(result, ", ")
}
