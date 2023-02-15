package huobi

import (
	"bytes"
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform"
)

func GetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var huobiJsonData = map[string]interface{}{
		"coinId":       token,     //Валюта 1-BTC, 2 - USDT
		"currency":     11,        //fiat 11-rub
		"tradeType":    tradeType, //купить - sell, продать - buy
		"currPage":     1,
		"payMethod":    platform.GetPayTypes(platform.HuobiName, c),
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

	result, err := platform.QueryToBytes(&huobiJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}
