package huobi

import "strings"

type Response struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	TotalCount int    `json:"totalCount"`
	PageSize   int    `json:"pageSize"`
	TotalPage  int    `json:"totalPage"`
	CurrPage   int    `json:"currPage"`
	Data       []struct {
		ID                int         `json:"id"`
		UserName          string      `json:"userName"`  //nickname
		CoinID            int         `json:"coinId"`    //token?
		Currency          int         `json:"currency"`  //fiat?
		TradeType         int         `json:"tradeType"` //1 -buy , 2-sell
		BlockType         int         `json:"blockType"`
		PayMethod         string      `json:"payMethod"`
		PayMethods        []PayMethod `json:"payMethods"`
		PayTerm           int         `json:"payTerm"`
		MinTradeLimit     string      `json:"minTradeLimit"`     //minLimit,
		MaxTradeLimit     string      `json:"maxTradeLimit"`     //maxLimit
		Price             string      `json:"price"`             //цена
		TradeCount        string      `json:"tradeCount"`        // доступно
		TradeMonthTimes   int         `json:"tradeMonthTimes"`   //Количество сделок
		OrderCompleteRate string      `json:"orderCompleteRate"` //процент выполнения
	} `json:"data"`
	Success bool `json:"success"`
}

type PayMethod struct {
	PayMethodID int    `json:"payMethodId"`
	Name        string `json:"name"`
}

func huobiPaymetodsToString(data []PayMethod) string {
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
