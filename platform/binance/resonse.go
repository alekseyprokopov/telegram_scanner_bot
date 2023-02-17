package binance

type Response struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          []struct {
		Adv struct {
			AdvNo                string `json:"advNo"`
			TradeType            string `json:"tradeType"`
			Asset                string `json:"asset"`
			FiatUnit             string `json:"fiatUnit"`
			Price                string `json:"price"`
			MaxSingleTransAmount string `json:"maxSingleTransAmount"`
			MinSingleTransAmount string `json:"minSingleTransAmount"`
			AutoReplyMsg         string `json:"autoReplyMsg"`
			TradeMethods         []struct {
				PayMethodID          string `json:"payMethodId"`
				Identifier           string `json:"identifier"`
				TradeMethodName      string `json:"tradeMethodName"`
				TradeMethodShortName string `json:"tradeMethodShortName"`
				TradeMethodBgColor   string `json:"tradeMethodBgColor"`
			} `json:"tradeMethods"`
			AssetScale                    int    `json:"assetScale"`
			FiatScale                     int    `json:"fiatScale"`
			PriceScale                    int    `json:"priceScale"`
			FiatSymbol                    string `json:"fiatSymbol"`
			IsTradable                    bool   `json:"isTradable"`
			DynamicMaxSingleTransAmount   string `json:"dynamicMaxSingleTransAmount"`
			MinSingleTransQuantity        string `json:"minSingleTransQuantity"`
			MaxSingleTransQuantity        string `json:"maxSingleTransQuantity"`
			DynamicMaxSingleTransQuantity string `json:"dynamicMaxSingleTransQuantity"`
			TradableQuantity              string `json:"tradableQuantity"`
			CommissionRate                string `json:"commissionRate"`
		} `json:"adv"`
		Advertiser struct {
			UserNo          string  `json:"userNo"`
			NickName        string  `json:"nickName"`
			MonthOrderCount int     `json:"monthOrderCount"`
			MonthFinishRate float64 `json:"monthFinishRate"`
			UserType        string  `json:"userType"`
			UserGrade       int     `json:"userGrade"`
			UserIdentity    string  `json:"userIdentity"`
		} `json:"advertiser"`
	} `json:"data"`
	Total   int  `json:"total"`
	Success bool `json:"success"`
}
