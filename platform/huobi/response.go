package huobi

import "strconv"

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

type SpotResponse struct {
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		Symbol  string  `json:"symbol"`
		Open    float64 `json:"open"`
		High    float64 `json:"high"`
		Low     float64 `json:"low"`
		Close   float64 `json:"close"`
		Amount  float64 `json:"amount"`
		Vol     float64 `json:"vol"`
		Count   int     `json:"count"`
		Bid     float64 `json:"bid"`
		BidSize float64 `json:"bidSize"`
		Ask     float64 `json:"ask"`
		AskSize float64 `json:"askSize"`
	} `json:"data"`
}

type PayMethod struct {
	PayMethodID int    `json:"payMethodId"`
	Name        string `json:"name"`
}

func getStringSlice(data []PayMethod) []string {
	var result []string
	for _, item := range data {
		result = append(result, strconv.Itoa(item.PayMethodID))
	}
	return result
}

func huobiTradeType(i int) string {
	if i == 1 {
		return "BUY"
	}
	return "SELL"
}
