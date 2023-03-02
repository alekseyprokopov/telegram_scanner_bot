package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scanner_bot/config"
	"strings"
	"sync"
)

type Platform interface {
	GetResult(c *config.Configuration) (*ResultPlatformData, error)
}

type PlatformTemplate struct {
	Name         string            `json:"platform_name"`
	Url          string            `json:"url"`
	ApiUrl       string            `json:"api_url"`
	Tokens       []string          `json:"platform_tokens"`
	TokensDict   map[string]string `json:"tokens_dict"`
	TradeTypes   []string          `json:"trade_types"`
	PayTypesDict map[string]string `json:"pay_types_dict"`
	AllPairs     map[string]bool   `json:"all_tokens"`
	Client       http.Client
}

func New(name string, url string, apiUrl string, tradeTypes []string, tokens []string, tokensDict map[string]string, payTypesDict map[string]string, allPairs map[string]bool) *PlatformTemplate {
	return &PlatformTemplate{
		Name:         name,
		ApiUrl:       apiUrl,
		Url:          url,
		TradeTypes:   tradeTypes,
		Tokens:       tokens,
		TokensDict:   tokensDict,
		PayTypesDict: payTypesDict,
		AllPairs:     allPairs,
		Client:       http.Client{},
	}
}

func (p *PlatformTemplate) TemplateResult(
	c *config.Configuration,
	spotData func() (*map[string]float64, error),
	advertise func(c *config.Configuration, token string, tradeType string) (*Advertise, error),
) (*ResultPlatformData, error) {
	result := ResultPlatformData{}
	result.Tokens = map[string]*TokenInfo{}
	wg := sync.WaitGroup{}
	var mu sync.Mutex
	result.Name = p.Name

	wg.Add(1)
	go func() {
		spotData, err := spotData()
		if err != nil {
			log.Printf("can't get spot data for %s: %v", p.Name, err)
		}
		mu.Lock()
		result.Spot = *spotData
		mu.Unlock()
		defer wg.Done()
	}()

	for _, token := range p.Tokens {
		token := token
		if p.Name == "huobi" {
			token = p.TokensDict[token]

		}
		tokenInfo := &TokenInfo{}
		result.Tokens[token] = tokenInfo

		wg.Add(1)
		go func() {
			buy, err := advertise(c, token, p.TradeTypes[0])
			if err != nil || buy == nil {
				log.Printf("can't get buy advertise for %s, token (%s): %v", p.Name, token, err)
			} else {
				mu.Lock()
				tokenInfo.Buy = *buy
				mu.Unlock()
			}
			defer wg.Done()
		}()

		wg.Add(1)
		go func() {
			sell, err := advertise(c, token, p.TradeTypes[1])
			if err != nil || sell == nil {
				log.Printf("can't get sell advertise for %s, token (%s): %v", p.Name, token, err)
			} else {
				mu.Lock()
				tokenInfo.Sell = *sell
				mu.Unlock()
			}
			defer wg.Done()
		}()
	}
	wg.Wait()

	return &result, nil

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

func (p *PlatformTemplate) DoGetRequest(urlAdd string, encodeQuery string) (*[]byte, error) {
	req, err := http.NewRequest(http.MethodGet, urlAdd, nil)
	if encodeQuery != "" {
		req.URL.RawQuery = encodeQuery
	}

	if err != nil {
		return nil, fmt.Errorf("can't do get request (%s): %w", p.Name, err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get resposnse from DoGetRequest (%s): %w", p.Name, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read info from response: %w", err)
	}

	return &body, err

}
func (p *PlatformTemplate) DoPostRequest(query *bytes.Buffer) (*[]byte, error) {
	req, err := http.NewRequest(http.MethodPost, p.Url, query)
	if err != nil {
		return nil, fmt.Errorf("can't create %s postRequest: %w", p.Name, err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("user-agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36`)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body: %w", err)
	}
	return &body, nil
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
