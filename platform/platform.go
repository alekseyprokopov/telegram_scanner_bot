package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"scanner_bot/config"
	"scanner_bot/platform/binance"
	"scanner_bot/platform/bybit"
	"scanner_bot/platform/huobi"
)

func New() *platformHandler {
	return &platformHandler{
		platformsInfo: map[string]Url{
			BinanceName: Url{
				P2PHost:   p2pBinanceHost,
				P2PPath:   p2pBinancePath,
				ApiUrl:    binanceAPI,
				Tokens:    binanceTokens,
				TradeType: binanceTradeType,
			},
			ByBitName: {
				P2PHost:   p2pByBitHost,
				P2PPath:   p2pByBitPath,
				ApiUrl:    byBitAPI,
				Tokens:    bybitTokens,
				TradeType: bybitTradeType,
			},
			HuobiName: {
				P2PHost:   p2pHuobiHost,
				P2PPath:   p2pHuobiPath,
				ApiUrl:    huobiAPI,
				Tokens:    huobiTokens,
				TradeType: huobiTradeType,
			},
		},
		client: http.Client{},
	}
}

func (p *platformHandler) TEST() {

}

func (p *platformHandler) GetAdvertise(platformName string, query *bytes.Buffer) ([]byte, error) {

	u := url.URL{
		Scheme: "https",
		Host:   p.platformsInfo[platformName].P2PHost,
		Path:   p.platformsInfo[platformName].P2PPath,
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

func getQuery(platformName string, c *config.Config, token string, tradeType string) (result *bytes.Buffer, err error) {
	switch platformName {
	case BinanceName:
		return binance.GetQuery(c, token, tradeType)
	case HuobiName:
		return huobi.GetQuery(c, token, tradeType)
	case ByBitName:
		return bybit.GetQuery(c, token, tradeType)
		//case GarantexName:
		//	return garantex.GetQuery(c, token, tradeType)
	}

	return nil, err
}

func GetPayTypes(platformName string, c *config.Config) []string {
	var result []string

	for key, value := range c.PayTypes {
		if value {
			result = append(result, payTypesDict[platformName][key])
		}
	}

	return result
}
