package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	BinanceURL  = "p2p.binance.com"
	BinancePath = "bapi/c2c/v2/friendly/c2c/adv/search"
)

type Platform struct {
	PlatformName     string          `json:"platform_name"`
	PayTypes         []string        `json:"pay_types"`
	IsActivePlatform bool            `json:"is_active_platform"`
	Roles            map[string]bool `json:"roles"`
	client           http.Client
}

func (p *Platform) GetData(host string, path string, buffer *bytes.Buffer) (map[string]interface{}, error) {

	resp, err := p.DoRequest(host, path, buffer)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result, nil

}

func (p *Platform) DoRequest(host string, path string, query *bytes.Buffer) ([]byte, error) {

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
	sellerName string
	bankName   string
	cost       int
	amount     int
	deals      int
}
