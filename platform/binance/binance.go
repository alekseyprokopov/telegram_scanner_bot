package binance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"io"
	"log"
	"net/http"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strconv"
	"strings"
	"sync"
)

type Platform struct {
	*platform.PlatformTemplate
	Binance *binance.Client
}

func New(name string, url string, tradeTypes []string, tokens []string, tokensDict map[string]string, payTypesDict map[string]string, allPairs map[string]bool) *Platform {
	return &Platform{
		PlatformTemplate: platform.New(name,url, tradeTypes, tokens, tokensDict, payTypesDict, allPairs),
		Binance: binance.NewClient("", ""),
	}
}

func (p *Platform) GetResult(c *config.Configuration) (*platform.ResultPlatformData, error) {
	result := platform.ResultPlatformData{}
	wg := sync.WaitGroup{}
	result.Name = p.Name

	wg.Add(1)
	go func() {
		spotData, err := p.getSpotData()
		if err != nil {
			log.Printf("can't get spot data: %v", err)
		}
		result.Spot = *spotData
		defer wg.Done()
	}()


	result.Tokens = map[string]*platform.TokenInfo{}

	for _, token := range p.Tokens {
		token:=token
		tokenInfo := &platform.TokenInfo{}
		result.Tokens[token] = tokenInfo

		wg.Add(1)
		go func() {
			buy, err := p.getAdvertise(c, token, p.TradeTypes[0])
			log.Println(token, buy.Cost)
			if err != nil || buy == nil {
				log.Printf("can't get buy advertise for huobi, token (%s): %v", token, err)
			} else {
				tokenInfo.Buy = *buy
			}
			defer wg.Done()


		}()


		wg.Add(1)
		go func() {
			sell, err := p.getAdvertise(c, token, p.TradeTypes[1])
			log.Println(token, sell.Cost)
			if err != nil || sell == nil {
				log.Printf("can't get sell advertise for huobi, token (%s): %v", token, err)
			} else {
				tokenInfo.Sell = *sell
			}
			defer wg.Done()
		}()
		//result.Tokens[token] = tokenInfo

	}
	wg.Wait()

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
		return nil, fmt.Errorf("can't do request to get binance response: %w", err)
	}

	advertise, err := p.responseToAdvertise(&response)
	if err != nil {
		return nil, fmt.Errorf("can't convert response to Advertise: %w", err)
	}

	return advertise, nil
}

func (p *Platform) getSpotData() (*map[string]float64, error) {
	set := p.AllPairs
	result := map[string]float64{}
	go p.test()
	prices, err := p.Binance.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't get spot data: %w", err)
	}
	for _, p := range prices {
		if set[p.Symbol] {
			result[p.Symbol], _ = strconv.ParseFloat(p.Price, 64)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("can't parse price to float: %w", err)
	}
	return &result, nil
}

func (p *Platform) getQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {

	var BinanceJsonData = map[string]interface{}{
		"proMerchantAds": false,
		"page":           1,
		"rows":           10,
		"payTypes":       p.GetPayTypes(c),
		"countries":      []string{},
		"publisherType":  nil,
		"transAmount":    c.MinValue,
		"asset":          token,
		"fiat":           "RUB",
		"tradeType":      tradeType, /*BUY(Купить) or SELL(продать)*/
	}

	result, err := p.QueryToBytes(&BinanceJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}

func (p *Platform) responseToAdvertise(response *[]byte) (*platform.Advertise, error) {
	var data Response
	err := json.Unmarshal(*response, &data)

	if err != nil {
		return nil, fmt.Errorf("cant' unmarshall data from binance response: %w", err)
	}
	item := data.Data[0]

	cost, _ := strconv.ParseFloat(item.Adv.Price, 64)
	minLimit, _ := strconv.ParseFloat(item.Adv.MinSingleTransAmount, 64)
	maxLimit, _ := strconv.ParseFloat(item.Adv.MaxSingleTransAmount, 64)
	available, _ := strconv.ParseFloat(item.Adv.DynamicMaxSingleTransQuantity, 64)
	return &platform.Advertise{
		PlatformName: p.Name,
		SellerName:   item.Advertiser.NickName,
		Asset:        item.Adv.Asset,
		Fiat:         item.Adv.FiatUnit,
		BankName:     payTypesToString(&data),
		Cost:         cost,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		SellerDeals:  item.Advertiser.MonthOrderCount,
		TradeType:    item.Adv.TradeType,
		Available:    available,
	}, nil
}

func payTypesToString(r *Response) string {
	data := r.Data[0].Adv.TradeMethods
	var result []string
	for _, k := range data {
		result = append(result, k.TradeMethodName)

	}
	return strings.Join(result, ", ")
}
func (p *Platform) doRequest(query *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, p.Url, query)
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

func (p *Platform) test()  {
	fmt.Println("GO ROUTINE")
}
