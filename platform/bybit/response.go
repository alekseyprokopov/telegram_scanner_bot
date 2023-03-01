package bybit

type Response struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  struct {
		Count int `json:"count"`
		Items []struct {
			NickName          string   `json:"nickName"`
			TokenName         string   `json:"tokenName"`    //tokenName
			TokenID           string   `json:"tokenId"`      //tokenName
			CurrencyID        string   `json:"currencyId"`   //фиат
			Side              int      `json:"side"`         //buy or sell
			Price             string   `json:"price"`        //цена
			LastQuantity      string   `json:"lastQuantity"` //доступно
			MinAmount         string   `json:"minAmount"`    //minLimit
			MaxAmount         string   `json:"maxAmount"`    //maxLimit
			Payments          []string `json:"payments"`
			RecentOrderNum    int      `json:"recentOrderNum"`    // количество сделок
			RecentExecuteRate int      `json:"recentExecuteRate"` //% выполнения
		} `json:"items"`
	} `json:"result"`
}


type SpotResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Category string `json:"category"`
		List     []struct {
			Symbol                 string `json:"symbol"`
			LastPrice              string `json:"lastPrice"`
			IndexPrice             string `json:"indexPrice"`
			MarkPrice              string `json:"markPrice"`
			PrevPrice24H           string `json:"prevPrice24h"`
			Price24HPcnt           string `json:"price24hPcnt"`
			HighPrice24H           string `json:"highPrice24h"`
			LowPrice24H            string `json:"lowPrice24h"`
			PrevPrice1H            string `json:"prevPrice1h"`
			OpenInterest           string `json:"openInterest"`
			OpenInterestValue      string `json:"openInterestValue"`
			Turnover24H            string `json:"turnover24h"`
			Volume24H              string `json:"volume24h"`
			FundingRate            string `json:"fundingRate"`
			NextFundingTime        string `json:"nextFundingTime"`
			PredictedDeliveryPrice string `json:"predictedDeliveryPrice"`
			BasisRate              string `json:"basisRate"`
			DeliveryFeeRate        string `json:"deliveryFeeRate"`
			DeliveryTime           string `json:"deliveryTime"`
			Ask1Size               string `json:"ask1Size"`
			Bid1Price              string `json:"bid1Price"`
			Ask1Price              string `json:"ask1Price"`
			Bid1Size               string `json:"bid1Size"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}