package platforms

import "net/http"

type PlatformHandler struct {
	platformsInfo map[string]Url
	client        http.Client
}

type Url struct {
	P2PHost   string
	P2PPath   string
	ApiUrl    string
	Tokens    []string
	TradeType []string
}

type Platform struct {
	PlatformName     string          `json:"platform_name"`
	PayTypes         []string        `json:"pay_types"`
	IsActivePlatform bool            `json:"is_active_platform"`
	Roles            map[string]bool `json:"roles"`
	client           http.Client
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
