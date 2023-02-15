package platform

import "net/http"

type platformHandler struct {
	platformsInfo map[string]Url
	client        http.Client
}

type Url struct {
	P2PHost string
	P2PPath string
	ApiUrl  string
	Tokens  []string
	TradeType []string
}

type Platform struct {
	PlatformName     string          `json:"platform_name"`
	PayTypes         []string        `json:"pay_types"`
	IsActivePlatform bool            `json:"is_active_platform"`
	Roles            map[string]bool `json:"roles"`
	client           http.Client
}

type ExchageInfo struct {
	name   string
	tokens map[string]tokenInfo
}

type tokenInfo struct {
	buy  Advertise
	sell Advertise
	spot []map[string]string
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
