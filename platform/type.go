package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"scanner_bot/config"
	"strings"
)

type Platform interface {
	GetResult(c *config.Configuration) (*ResultPlatformData, error)

}

type PlatformTemplate struct {
	Name         string            `json:"platform_name"`
	Url          string            `json:"url"`
	Tokens       []string          `json:"platform_tokens"`
	TokensDict   map[string]string `json:"tokens_dict"`
	TradeTypes   []string          `json:"trade_types"`
	PayTypesDict map[string]string `json:"pay_types_dict"`
	AllPairs     map[string]bool   `json:"all_tokens"`
	Client       http.Client
}

func New(name string, url string, tradeTypes []string, tokens []string, tokensDict map[string]string, payTypesDict map[string]string, allPairs map[string]bool) *PlatformTemplate {
	return &PlatformTemplate{
		Name:         name,
		Url:          url,
		TradeTypes:   tradeTypes,
		Tokens:       tokens,
		TokensDict:   tokensDict,
		PayTypesDict: payTypesDict,
		AllPairs:     allPairs,
		Client:       http.Client{},
	}
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
			result = append(result, p.PayTypesDict[key])
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

func (p *PlatformTemplate) PayTypesToString(data []string) string {
	var result []string
	for _, item := range data {
		item, ok := p.PayTypesDict[item]
		if ok {
			result = append(result, item)
		}
	}
	return strings.Join(result, ", ")
}

func (p *PlatformTemplate) TokenFromDict(item string) string {
	return p.TokensDict[item]
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
	Tokens map[string]*TokenInfo
}

type TokenInfo struct {
	Buy  Advertise
	Sell Advertise
}
