package bybit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hirokisan/bybit/v2"
	"io"
	"net/http"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
)

type Platform struct {
	*platform.PlatformTemplate
	Bybit *bybit.Client
}

func New(name string, url string, tradeTypes []string, tokens []string, tokensDict map[string]string, payTypesDict map[string]string, allPairs map[string]bool) *Platform {
	return &Platform{
		PlatformTemplate: platform.New(name, url, tradeTypes, tokens, tokensDict, payTypesDict, allPairs),
		Bybit:            bybit.NewClient().WithAuth("", ""),
	}
}

func (p *Platform) GetResult(c *config.Configuration) (*platform.ResultPlatformData, error) {
	result := platform.ResultPlatformData{}
	spotData, err := p.getSpotData()
	if err != nil {
		return nil, fmt.Errorf("can't get bybit spot data: %w", err)
	}

	result.Name = p.Name
	result.Spot = *spotData
	result.Tokens = map[string]platform.TokenInfo{}
	if err != nil {
		return nil, fmt.Errorf("cant get advertise: %w", err)
	}
	for _, token := range p.Tokens {
		buy, err := p.getAdvertise(c, token, p.TradeTypes[0])
		sell, err := p.getAdvertise(c, token, p.TradeTypes[1])
		if err != nil {
			return nil, fmt.Errorf("can't get advertise: %w", err)
		}
		result.Tokens[token] = platform.TokenInfo{
			Buy:  *buy,
			Sell: *sell,
		}
	}
	return &result, nil

}

func (p *Platform) getSpotData() (*map[string]float64, error) {

	set := p.AllPairs
	result := map[string]float64{}

	sym := bybit.SymbolFuture("")
	response, _ := p.Bybit.Future().InverseFuture().Tickers(sym)
	for _, item := range response.Result {
		_, ok := set[string(item.Symbol)]
		if ok {
			result[string(item.Symbol)], _ = strconv.ParseFloat(item.LastPrice, 64)
		}
	}
	return &result, nil
}

func (p *Platform) getAdvertise(c *config.Configuration, token string, tradeType string) (*platform.Advertise, error) {
	userConfig := &c.UserConfig

	query, err := p.getQuery(userConfig, token, tradeType)
	if err != nil {
		return nil, fmt.Errorf("can't get query: %w", err)
	}
	response, err := p.doRequest(query)
	if err != nil {
		return nil, fmt.Errorf("can't do request to get bybit response: %w", err)
	}

	advertise, err := p.responseToAdvertise(&response)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}

	return advertise, nil
}

func (p *Platform) getQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var BybitJsonData = map[string]interface{}{
		"userId":     "",
		"tokenId":    token,
		"currencyId": "RUB",
		"payment":    p.GetPayTypes(c),
		"side":       tradeType,
		"size":       "1",
		"page":       "1",
		"amount":     strconv.Itoa(c.MinValue),
	}
	result, err := p.QueryToBytes(&BybitJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

func (p *Platform) doRequest(query *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, p.Url, query)
	if err != nil {
		return nil, fmt.Errorf("can't do bybit request: %w", err)
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
	return body, nil
}

func (p *Platform) responseToAdvertise(response *[]byte) (*platform.Advertise, error) {
	var data Response
	err := json.Unmarshal(*response, &data)
	if err != nil {
		return nil, fmt.Errorf("cant' unmarshall data from bybit response: %w", err)
	}
	item := data.Result.Items[0]

	cost, _ := strconv.ParseFloat(item.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.MinAmount, 64)
	maxLimit, _ := strconv.ParseFloat(item.MaxAmount, 64)
	available, _ := strconv.ParseFloat(item.LastQuantity, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.NickName,
		Asset:        item.TokenName,
		Fiat:         item.CurrencyID,
		BankName:     p.PayTypesToString(item.Payments),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.RecentOrderNum,
		TradeType:    bybitTradeType(item.Side),
		Available:    available,
	}, nil
}

func bybitTradeType(i int) string {
	if i == 1 {
		return "BUY"
	}
	return "SELL"
}

// сделать функцию универсальной!!!
