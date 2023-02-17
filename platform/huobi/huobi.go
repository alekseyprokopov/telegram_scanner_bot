package huobi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
)

type Platform struct {
	platform.PlatformTemplate
}

func (p *Platform) GetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var huobiJsonData = map[string]interface{}{
		"coinId":       token,     //Валюта 1-BTC, 2 - USDT
		"currency":     11,        //fiat 11-rub
		"tradeType":    tradeType, //купить - sell, продать - buy
		"currPage":     1,
		"payMethod":    p.GetPayTypes(c),
		"acceptOrder":  0,
		"country":      "",
		"blockType":    "general",
		"online":       "1",
		"range":        0,
		"amount":       c.MinValue,
		"isThumbsUp":   false,
		"isMerchant":   false,
		"isTraded":     false,
		"onlyTradable": false,
		"isFollowed":   false,
	}

	result, err := p.QueryToBytes(&huobiJsonData)
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

	cost, _ := strconv.ParseFloat(item.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.MinTradeLimit, 64)
	maxLimit, _ := strconv.ParseFloat(item.MaxTradeLimit, 64)
	available, _ := strconv.ParseFloat(item.TradeCount, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.UserName,
		Asset:        string(item.CoinID),
		Fiat:         string(item.Currency),
		BankName:     huobiPaymetodsToString(item.PayMethods),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.TradeMonthTimes,
		TradeType:    huobiTradeType(item.TradeType),
		Available:    available,
	}, nil
}


