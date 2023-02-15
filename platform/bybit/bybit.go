package bybit

import (
	"bytes"
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform"
)

func GetQuery(c *config.Config, token string, tradeType string) (*bytes.Buffer, error) {
	var byBitJsonData = map[string]interface{}{
		"userId":     "",
		"tokenId":    token,
		"currencyID": "RUB",
		"payment":    platform.GetPayTypes(platform.ByBitName, c),
		"side":       "1",
		"size":       "10",
		"page":       "1",
		"amount":     c.MinValue,
	}

	result, err := platform.QueryToBytes(&byBitJsonData)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return result, nil
}
