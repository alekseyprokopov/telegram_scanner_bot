package bybit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
	"strings"
)

type Platform struct {
	platform.PlatformTemplate
}

func (p *Platform) GetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var byBitJsonData = map[string]interface{}{
		"userId":     "",
		"tokenId":    token,
		"currencyID": "RUB",
		"payment":    p.GetPayTypes(c),
		"side":       tradeType,
		"size":       "10",
		"page":       "1",
		"amount":     c.MinValue,
	}

	result, err := p.QueryToBytes(&byBitJsonData)
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
	item := data.Result.Items[0]

	cost, _ := strconv.ParseFloat(item.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.MinAmount, 64)
	maxLimit, _ := strconv.ParseFloat(item.MaxAmount, 64)
	available, _ := strconv.ParseFloat(item.LastQuantity, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.NickName,
		Asset:        item.TokenName,
		Fiat:         item.CurrencyID,
		BankName:     bybitPayTypesToString(item.Payments),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.RecentOrderNum,
		TradeType:    bybitTradeType(item.Side),
		Available:    available,
	}, nil
}

func bybitTradeType(i int) string {
	if i == 1 {
		return "BUY"
	}
	return "SELL"
}

func bybitPayTypesToString(item []int) string {
	dict := platform.PayTypesDict[platform.ByBitName]
	var result []string

	for _, value := range item {
		result = append(result, dict[string(value)])
	}
	return strings.Join(result, ", ")
}
