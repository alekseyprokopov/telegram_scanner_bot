package platforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
	"strconv"
	"strings"
)

func huobiGetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var huobiJsonData = map[string]interface{}{
		"coinId":       token,     //Валюта 1-BTC, 2 - USDT
		"currency":     11,        //fiat 11-rub
		"tradeType":    tradeType, //купить - sell, продать - buy
		"currPage":     1,
		"payMethod":    GetPayTypes(HuobiName, c),
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

	result, err := QueryToBytes(&huobiJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

type huobiResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	TotalCount int    `json:"totalCount"`
	PageSize   int    `json:"pageSize"`
	TotalPage  int    `json:"totalPage"`
	CurrPage   int    `json:"currPage"`
	Data       []struct {
		ID                int              `json:"id"`
		UserName          string           `json:"userName"`  //nickname
		CoinID            int              `json:"coinId"`    //token?
		Currency          int              `json:"currency"`  //fiat?
		TradeType         int              `json:"tradeType"` //1 -buy , 2-sell
		BlockType         int              `json:"blockType"`
		PayMethod         string           `json:"payMethod"`
		PayMethods        []huobiPayMethod `json:"payMethods"`
		PayTerm           int              `json:"payTerm"`
		MinTradeLimit     string           `json:"minTradeLimit"`     //minLimit,
		MaxTradeLimit     string           `json:"maxTradeLimit"`     //maxLimit
		Price             string           `json:"price"`             //цена
		TradeCount        string           `json:"tradeCount"`        // доступно
		TradeMonthTimes   int              `json:"tradeMonthTimes"`   //Количество сделок
		OrderCompleteRate string           `json:"orderCompleteRate"` //процент выполнения
	} `json:"data"`
	Success bool `json:"success"`
}

type huobiPayMethod struct {
	PayMethodID int    `json:"payMethodId"`
	Name        string `json:"name"`
}

func HuobiResponseToAdvertise(response *[]byte) (*Advertise, error) {
	var data huobiResponse
	err := json.Unmarshal(*response, &data)
	if err != nil {
		return nil, fmt.Errorf("cant' unmarshall data from binance response: %w", err)
	}
	item := data.Data[0]

	cost, _ := strconv.ParseFloat(item.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.MinTradeLimit, 64)
	maxLimit, _ := strconv.ParseFloat(item.MaxTradeLimit, 64)
	available, _ := strconv.ParseFloat(item.TradeCount, 64)
	return &Advertise{
		PlatformName: BinanceName,
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

func huobiPaymetodsToString(data []huobiPayMethod) string {
	var result []string

	for _, item := range data {
		result = append(result, item.Name)
	}

	return strings.Join(result, ", ")
}


func huobiTradeType(i int) string {
	if i == 1 {
		return "BUY"
	}
	return "SELL"
}