package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scanner_bot/config"
)

type PlatformHandler struct {
	platformsInfo map[string]Platform
	client        http.Client
}

type Platform interface {
	//GetAdvertise(c *config.Configuration, token string, tradeType string) (*Advertise, error)
	GetResult(c *config.Configuration) (*ResultPlatformData, error)
	//GetQuery()
	//ResponseToAdvertise()
}

type PlatformTemplate struct {
	Name       string   `json:"platform_name"`
	Url        string   `json:"url"`
	PayTypes   []string `json:"pay_types"`
	TradeTypes []string `json:"trade_types"`
	Tokens     []string `json:"platform_tokens"`
	AllTokens  []string `json:"all_tokents"`
	Client     http.Client
}

func (p *PlatformTemplate) QueryToBytes(params *map[string]interface{}) (*bytes.Buffer, error) {
	bytesRepresentation, err := json.Marshal(*params)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return bytes.NewBuffer(bytesRepresentation), nil
}

func (p *PlatformTemplate) GetPayTypes(c *config.Config) []string {
	var result []string
	for key, value := range c.PayTypes {
		if value {
			result = append(result, PayTypesDict[p.Name][key])
		}
	}
	return result
}
func (p *PlatformTemplate) CreatePairsSet(data []string) *map[string]bool {
	set := map[string]bool{}
	for i := 0; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			set[data[i]+data[j]] = true
			set[data[j]+data[i]] = true
		}
	}
	return &set
}

func (p *PlatformTemplate) DoRequest(query *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, p.Url, query)
	if err != nil {
		return nil, fmt.Errorf("can't do binance request: %w", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", `"Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"`)

	resp, err := p.Client.Do(req)
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
	Name   string
	Spot   map[string]float64
	Tokens map[string]TokenInfo
}

type TokenInfo struct {
	Buy  Advertise
	Sell Advertise
}
