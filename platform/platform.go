package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	binanceTokens = []string{"USTD", "BTC", "BUSD", "BNB", "ETH", "SHIB"}
	bybitTokens   = []string{"USDT", "BTC", "ETH", "USDC"}
	huobiTokens   = []string{
		"2",  /*USDT*/
		"1",  /*BTC*/
		"62", /*USDD*/
		"4",  /*HT*/
		"22", /*TRX*/
		"3",  /*ETH*/
		"7",  /*XRP*/
		"8",  /*LTC*/
	}
)

const (
	binanceName    = "binance"
	p2pBinanceHost = "p2p.binance.com"
	p2pBinancePath = "bapi/c2c/v2/friendly/c2c/adv/search"
	binanceAPI     = ""

	byBitName    = "bybit"
	p2pByBitHost = "api2.bybit.com"
	p2pByBitPath = "spot/api/otc/item/list"
	byBitAPI     = ""

	huobiName    = "huobi"
	p2pHuobiHost = "otc-akm.huobi.com"
	p2pHuobiPath = "v1/data/trade-market?"
	huobiAPI     = ""
)

type platformHandler map[string]Url

type Url struct {
	P2PHost string
	P2PPath string
	ApiUrl  string
	Tokens  []string
}

func New() *platformHandler {
	return &platformHandler{
		binanceName: Url{
			P2PHost: p2pBinanceHost,
			P2PPath: p2pBinancePath,
			ApiUrl:  binanceAPI,
			Tokens:  binanceTokens,
		},
		byBitName: {
			P2PHost: p2pByBitHost,
			P2PPath: p2pByBitPath,
			ApiUrl:  byBitAPI,
			Tokens:  bybitTokens,
		},
		huobiName: {
			P2PHost: p2pHuobiHost,
			P2PPath: p2pHuobiPath,
			ApiUrl:  huobiAPI,
			Tokens:  huobiTokens,
		},
	}
}

type Platform struct {
	PlatformName     string          `json:"platform_name"`
	PayTypes         []string        `json:"pay_types"`
	IsActivePlatform bool            `json:"is_active_platform"`
	Roles            map[string]bool `json:"roles"`
	client           http.Client
}

func (p *Platform) GetData(host string, path string, query *bytes.Buffer) ([]byte, error) {

	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), query)
	if err != nil {
		return nil, fmt.Errorf("can't do binance request: %w", err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", `"Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"`)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body: %w", err)
	}
	return body, nil

}

func QueryToBytes(params *map[string]interface{}) (*bytes.Buffer, error) {
	bytesRepresentation, err := json.Marshal(*params)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return bytes.NewBuffer(bytesRepresentation), nil
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
