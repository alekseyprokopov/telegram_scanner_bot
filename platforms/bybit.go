package platforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
	"strconv"
	"strings"
)

func bybitGetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var byBitJsonData = map[string]interface{}{
		"userId":     "",
		"tokenId":    token,
		"currencyID": "RUB",
		"payment":    GetPayTypes(ByBitName, c),
		"side":       "1",
		"size":       "10",
		"page":       "1",
		"amount":     c.MinValue,
	}

	result, err := QueryToBytes(&byBitJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

type bybitResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  struct {
		Count int `json:"count"`
		Items []struct {
			NickName          string `json:"nickName"`
			TokenName         string `json:"tokenName"`    //tokenName
			CurrencyID        string `json:"currencyId"`   //фиат
			Side              int    `json:"side"`         //buy or sell
			Price             string `json:"price"`        //цена
			LastQuantity      string `json:"lastQuantity"` //доступно
			MinAmount         string `json:"minAmount"`    //minLimit
			MaxAmount         string `json:"maxAmount"`    //maxLimit
			Payments          []int  `json:"payments"`
			RecentOrderNum    int    `json:"recentOrderNum"`    // количество сделок
			RecentExecuteRate int    `json:"recentExecuteRate"` //% выполнения
		} `json:"items"`
	} `json:"result"`
	Token   interface{} `json:"token"`
	ExtInfo interface{} `json:"ext_info"`
}

func BybitResponseToAdvertise(response *[]byte) (*Advertise, error) {
	var data bybitResponse
	item := data.Result.Items[0]
	err := json.Unmarshal(*response, &data)
	if err != nil {
		return nil, fmt.Errorf("cant' unmarshall data from binance response: %w", err)
	}

	cost, _ := strconv.ParseFloat(item.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.MinAmount, 64)
	maxLimit, _ := strconv.ParseFloat(item.MaxAmount, 64)
	available, _ := strconv.ParseFloat(item.LastQuantity, 64)
	return &Advertise{
		PlatformName: BinanceName,
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
	dict := PayTypesDict[ByBitName]
	var result []string

	for _, value := range item {
		result = append(result, dict[string(value)])
	}
	return strings.Join(result, ", ")
}
