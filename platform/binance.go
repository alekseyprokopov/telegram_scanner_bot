package platform

import (
	"bytes"
	"fmt"
	"scanner_bot/config"
	"strconv"
	"strings"
)

func binanceGetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {

	var BinanceJsonData = map[string]interface{}{
		"proMerchantAds": false,
		"page":           1,
		"rows":           10,
		"payTypes":       GetPayTypes(BinanceName, c),
		"countries":      []string{},
		"publisherType":  nil,
		"transAmount":    c.MinValue,
		"asset":          token,
		"fiat":           "RUB",
		"tradeType":      tradeType, /*BUY(Купить) or SELL(продать)*/
	}

	result, err := QueryToBytes(&BinanceJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

type binanceResponse struct {
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

func binanceResponseToAdvertise(r *binanceResponse) *Advertise {
	cost, _ := strconv.ParseFloat(r.Data[0].Adv.Price, 64)
	minLimit, _ := strconv.ParseFloat(r.Data[0].Adv.MinSingleTransAmount, 64)
	maxLimit, _ := strconv.ParseFloat(r.Data[0].Adv.MaxSingleTransAmount, 64)
	available, _ := strconv.ParseFloat(r.Data[0].Adv.DynamicMaxSingleTransQuantity, 64)
	return &Advertise{
		PlatformName: BinanceName,
		SellerName:   r.Data[0].Advertiser.NickName,
		Asset:        r.Data[0].Adv.Asset,
		Fiat:         r.Data[0].Adv.FiatUnit,
		BankName:     binancePayTypesToString(r),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  r.Data[0].Advertiser.MonthOrderCount,
		TradeType:    r.Data[0].Adv.TradeType,
		Available:    available,
	}
}

func binancePayTypesToString(r *binanceResponse) string {
	data := r.Data[0].Adv.TradeMethods
	var result []string
	for _, k := range data {
		result = append(result, k.TradeMethodName)

	}
	return strings.Join(result, ", ")
}
