package binance

type Response struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          []AdvertiseData `json:"data"`
	Total   int  `json:"total"`
	Success bool `json:"success"`
}

type AdvertiseData struct {
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
}

type SpotResponse struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		BaseAsset                  string   `json:"baseAsset"`
		BaseAssetPrecision         int      `json:"baseAssetPrecision"`
		QuoteAsset                 string   `json:"quoteAsset"`
		QuotePrecision             int      `json:"quotePrecision"`
		QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
		BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
		QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
		OrderTypes                 []string `json:"orderTypes"`
		IcebergAllowed             bool     `json:"icebergAllowed"`
		OcoAllowed                 bool     `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		AllowTrailingStop          bool     `json:"allowTrailingStop"`
		CancelReplaceAllowed       bool     `json:"cancelReplaceAllowed"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType            string `json:"filterType"`
			MinPrice              string `json:"minPrice,omitempty"`
			MaxPrice              string `json:"maxPrice,omitempty"`
			TickSize              string `json:"tickSize,omitempty"`
			MinQty                string `json:"minQty,omitempty"`
			MaxQty                string `json:"maxQty,omitempty"`
			StepSize              string `json:"stepSize,omitempty"`
			MinNotional           string `json:"minNotional,omitempty"`
			ApplyToMarket         bool   `json:"applyToMarket,omitempty"`
			AvgPriceMins          int    `json:"avgPriceMins,omitempty"`
			Limit                 int    `json:"limit,omitempty"`
			MinTrailingAboveDelta int    `json:"minTrailingAboveDelta,omitempty"`
			MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta,omitempty"`
			MinTrailingBelowDelta int    `json:"minTrailingBelowDelta,omitempty"`
			MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta,omitempty"`
			BidMultiplierUp       string `json:"bidMultiplierUp,omitempty"`
			BidMultiplierDown     string `json:"bidMultiplierDown,omitempty"`
			AskMultiplierUp       string `json:"askMultiplierUp,omitempty"`
			AskMultiplierDown     string `json:"askMultiplierDown,omitempty"`
			MaxNumOrders          int    `json:"maxNumOrders,omitempty"`
			MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions                     []string `json:"permissions"`
		DefaultSelfTradePreventionMode  string   `json:"defaultSelfTradePreventionMode"`
		AllowedSelfTradePreventionModes []string `json:"allowedSelfTradePreventionModes"`
	} `json:"symbols"`
}
