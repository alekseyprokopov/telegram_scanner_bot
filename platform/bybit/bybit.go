package bybit

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
		//Binance: binance.NewClient("", ""),
	}
}

func (p *Platform) GetResult(c *config.Configuration) (*platform.ResultPlatformData, error) {
	adv, err := p.GetAdvertise(c, p.Tokens[0], p.TradeTypes[0])
	if err != nil {
		return nil, fmt.Errorf("can't get advertise: %w", err)
	}
	log.Printf("advertise : %+v", adv)
	return nil, nil
}
func (p *Platform) GetAdvertise(c *config.Configuration, token string, tradeType string) (*platform.Advertise, error) {
	userConfig := &c.UserConfig

	query := p.getQuery(userConfig, token, tradeType)

	log.Printf("query from bybit: %+v", query)
	response, err := p.doRequest(query)
	if err != nil {
		return nil, fmt.Errorf("can't do request to get bybit response: %w", err)
	}
	log.Printf("response from bybit: %+v", string(response))

	advertise, err := p.responseToAdvertise(&response)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}

	log.Printf("advertise: %+v", advertise)

	return advertise, nil
}

func (p *Platform) getQuery(c *config.Config, token string, tradeType string) string {
	q := url.Values{
		"userId":     []string{},
		"tokenId":    []string{token},
		"currencyId": []string{"RUB"},
		"payment":    p.GetPayTypes(c),
		"side":       []string{tradeType},
		"size":       []string{"1"}, // количество приходящих объявлений
		"page":       []string{"1"},
		"amount":     []string{strconv.Itoa(c.MinValue)},
	}
	result := q.Encode()
	return result
}

func (p *Platform) doRequest(query string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, p.Url, nil)
	req.URL.RawQuery = query

	if err != nil {
		return nil, fmt.Errorf("can't do binance request: %w", err)
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
	log.Print(item.Payments)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.NickName,
		Asset:        item.TokenName,
		Fiat:         item.CurrencyID,
		BankName:     PayTypesToString(item.Payments),
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
func PayTypesToString(item []int) string {
	dict := platform.PayTypesDict[platform.ByBitName]
	var result []string

	for _, value := range item {
		item, ok := dict[strconv.Itoa(value)]
		if ok {
			result = append(result, item)

		}
	}
	return strings.Join(result, ", ")
}
