package huobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
	"strings"
)

type Platform struct {
	*platform.PlatformTemplate
	//Huobi
}

func New(name string, url string, tradeTypes []string, payTypes []string, tokens []string, allTokens []string) *Platform {
	return &Platform{
		PlatformTemplate: &platform.PlatformTemplate{
			Name:       name,
			Url:        url,
			TradeTypes: tradeTypes,
			Tokens:     tokens,
			PayTypes:   payTypes,
			AllTokens:  allTokens,
			Client:     http.Client{},
		},
		//Huobi: bybit.NewClient().WithAuth("", ""),
	}
}

func (p *Platform) GetResult(c *config.Configuration) (*platform.ResultPlatformData, error) {
	result := platform.ResultPlatformData{}
	//spotData, err := p.getSpotData()
	//if err != nil {
	//	return nil, fmt.Errorf("can't get huobi spot data: %w", err)
	//}
	result.Name = p.Name
	//result.Spot = *spotData
	result.Tokens = map[string]platform.TokenInfo{}

	for _, token := range p.Tokens {
		tokenResult := platform.TokenInfo{}
		buy, err := p.getAdvertise(c, token, p.TradeTypes[0])
		if err != nil || buy == nil {
			log.Printf("can't get buy advertise for huobi, token (%s): %v", token, err)
		} else {
			tokenResult.Buy = *buy
		}
		sell, err := p.getAdvertise(c, token, p.TradeTypes[1])

		if err != nil || sell == nil {
			log.Printf("can't get sell advertise for huobi, token (%s): %v", token, err)
		} else {
			tokenResult.Sell = *sell
		}
		result.Tokens[token] = tokenResult
	}
	return &result, nil
}
func (p *Platform) getAdvertise(c *config.Configuration, token string, tradeType string) (*platform.Advertise, error) {
	userConfig := &c.UserConfig
	query := p.getQuery(userConfig, token, tradeType)
	response, err := p.doRequest(query)
	if err != nil {
		return nil, fmt.Errorf("can't do request to get bybit response: %w", err)
	}

	advertise, err := p.responseToAdvertise(response)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}

	return advertise, nil
}

func (p *Platform) getQuery(c *config.Config, token string, tradeType string) string {
	u := url.Values{
		"coinId":       []string{token},     //usdt
		"currency":     []string{"11"},      //rub
		"tradeType":    []string{tradeType}, //buy
		"currPage":     []string{"1"},
		"payMethod":    []string{strings.Join(p.GetPayTypes(c), ",")},
		"acceptOrder":  []string{"0"},
		"country":      []string{""},
		"blockType":    []string{"general"},
		"online":       []string{"1"},
		"range":        []string{"0"},
		"amount":       []string{strconv.Itoa(c.MinValue)}, //amount
		"isThumbsUp":   []string{"false"},
		"isMerchant":   []string{"false"},
		"isTraded":     []string{"false"},
		"onlyTradable": []string{"false"},
		"isFollowed":   []string{"false"},
	}

	return u.Encode()
}

func (p *Platform) doRequest(query string) (*[]byte, error) {
	req, err := http.NewRequest(http.MethodGet, p.Url, nil)
	req.URL.RawQuery = query
	if err != nil {
		return nil, fmt.Errorf("can't do huobi request: %w", err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/json")
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

func (p *Platform) responseToAdvertise(response *[]byte) (*platform.Advertise, error) {
	var data Response
	err := json.Unmarshal(*response, &data)
	if err != nil || len(data.Data) == 0 || data.Code != 200 {
		return nil, fmt.Errorf("cant' unmarshall data from binance response: %w", err)
	}
	item := data.Data[0]

	cost, _ := strconv.ParseFloat(item.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.MinTradeLimit, 64)
	maxLimit, _ := strconv.ParseFloat(item.MaxTradeLimit, 64)
	available, _ := strconv.ParseFloat(item.TradeCount, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.UserName,
		Asset:        huobiTokensFromDict(item.CoinID),
		Fiat:         strconv.Itoa(item.Currency) + " (RUB)",
		BankName:     payMethodsToString(item.PayMethods),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.TradeMonthTimes,
		TradeType:    huobiTradeType(item.TradeType),
		Available:    available,
	}, nil
}
