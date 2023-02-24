package bybit

type Response struct {
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
			Payments          []string  `json:"payments"`
			RecentOrderNum    int    `json:"recentOrderNum"`    // количество сделок
			RecentExecuteRate int    `json:"recentExecuteRate"` //% выполнения
		} `json:"items"`
	} `json:"result"`
}