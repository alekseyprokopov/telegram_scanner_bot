package platforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"scanner_bot/config"
)

func New() *PlatformHandler {
	return &PlatformHandler{
		platformsInfo: map[string]Url{
			BinanceName: Url{
				P2PHost:   P2PBinanceHost,
				P2PPath:   P2PBinancePath,
				ApiUrl:    BinanceAPI,
				Tokens:    BinanceTokens,
				TradeType: BinanceTradeType,
			},
			ByBitName: {
				P2PHost:   P2PByBitHost,
				P2PPath:   P2PByBitPath,
				ApiUrl:    ByBitAPI,
				Tokens:    BybitTokens,
				TradeType: BybitTradeType,
			},
			HuobiName: {
				P2PHost:   P2PHuobiHost,
				P2PPath:   P2PHuobiPath,
				ApiUrl:    HuobiAPI,
				Tokens:    HuobiTokens,
				TradeType: HuobiTradeType,
			},
		},
		client: http.Client{},
	}
}

func (p *PlatformHandler) GetAdvertise(c *config.Configuration) (*Advertise, error) {

	userConfig := &c.UserConfig
	platformName := BinanceName

	query, err := GetQuery(platformName, userConfig, "USDT", "BUY")

	if err != nil {
		return nil, fmt.Errorf("can't get query: %w", err)
	}


	response, err := p.getAdvertiseData(platformName, query)

	var Binance BinanceResponse

	json.Unmarshal(data, &Binance)

	var resultAdvertise = BinanceResponseToAdvertise(&Binance)
	log.Printf("advertise: %+v", resultAdvertise)
	result := msgAdvertise(info)
	log.Printf("result: %+v", result)

	return nil, resultAdvertise
}

func (p *PlatformHandler) getAdvertiseData(platformName string, query *bytes.Buffer) ([]byte, error) {
	log.Println("platformName :", platformName)
	log.Println("p.platformsInfo[platformName] :", p.platformsInfo[platformName])

	u := url.URL{
		Scheme: "https",
		Host:   p.platformsInfo[platformName].P2PHost,
		Path:   p.platformsInfo[platformName].P2PPath,
	}

	log.Println("url :", u.String())

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

func ResponseToAdvertise(platformName string, response *[]byte) (*Advertise, error) {
	switch platformName {
	case BinanceName:
		return BinanceResponseToAdvertise(response)
	case HuobiName:
		return HuobiRespinseToAdvertise(response)
	case ByBitName:
		return ByBitRespinseToAdvertise(response)


	}
}
