package platform

import "net/http"

type PlatformHandler struct {
	platformsInfo map[string]Platform
	client        http.Client
}

type Platform interface {
	GetAdvertise()
	GetResultInfo()
	GetQuery()
	ResponseToADvertise()
}

type PlatformType struct {
	Name             string   `json:"platform_name"`
	Url              string   `json:"url"`
	PayTypes         []string `json:"pay_types"`
	PlatformTokens   []string `json:"platform_tokens"`
	AllTokents       []string `json:"all_tokents"`
	Client           http.Client
}

type Advertise struct {
	PlatformName string
	SellerName   string
	Asset        string
	Fiat         string
	BankName     string
	Cost         float64
	MinLimit     float64
	MaxLimit     float64
	SellerDeals  int
	TradeType    string
	Available    float64
}

type ResultPlatformData struct {
	PlatformName string
	Tokens       map[string]TokenInfo
}

type TokenInfo struct {
	buy  Advertise
	sell Advertise
	spot []map[string]int
}
