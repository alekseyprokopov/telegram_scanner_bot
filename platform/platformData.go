package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform/binance"
	"scanner_bot/platform/bybit"
	"scanner_bot/platform/huobi"
)

var (
	BinanceTokens    = []string{"USTD", "BTC", "BUSD", "BNB", "ETH", "SHIB"}
	BinanceTradeType = []string{"BUY", "SELL"}

	BybitTokens    = []string{"USDT", "BTC", "ETH", "USDC"}
	BybitTradeType = []string{"1", "0"} //side 1- купить. 0 - продать

	HuobiTokens = []string{
		"2",  /*USDT*/
		"1",  /*BTC*/
		"62", /*USDD*/
		"4",  /*HT*/
		"22", /*TRX*/
		"3",  /*ETH*/
		"7",  /*XRP*/
		"8",  /*LTC*/
	}
	HuobiTradeType = []string{"sell", "buy"} //sell-купить, buy - продать

	PayTypesDict = map[string]map[string]string{
		BinanceName: {
			"Сбербанк": "RosBankNew", "Тинькофф": "TinkoffNew", "Райффайзен": "RaiffeisenBank", "QIWI": "QIWI", "ЮMoney": "YandexMoneyNew",
		},
		HuobiName: {
			"Сбербанк": "29", "Тинькофф": "28", "Райффайзен": "36", "QIWI": "9", "ЮMoney": "19",
		},
		ByBitName: {
			"185": "Сбербанк", "75": "Тинькофф", "64": "Райффайзен", "62": "QIWI", "274": "ЮMoney",
		},
		GarantexName: {
			"Сбербанк": "Сбербанк", "Тинькофф": "Тинькофф", "Райффайзен": "Райффайзен", "QIWI": "QIWI", "ЮMoney": "ЮMoney",
		},
	}
)

const (
	BinanceName   = "binance"
	P2PBinanceURL = "p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"

	ByBitName    = "bybit"
	P2PByBitURL  = "api2.bybit.com/spot/api/otc/item/list"


	HuobiName    = "huobi"
	P2PHuobiURL = "otc-akm.huobi.com/v1/data/trade-market?"

	GarantexName = "garantex"
)

func GetQuery(platformName string, c *config.Config, token string, tradeType string) (result *bytes.Buffer, err error) {
	switch platformName {
	case BinanceName:
		return binance.binanceGetQuery(c, token, tradeType)
	case HuobiName:
		return huobi.huobiGetQuery(c, token, tradeType)
	case ByBitName:
		return bybit.bybitGetQuery(c, token, tradeType)
		//case GarantexName:
		//	return garantex.GetQuery(c, token, tradeType)
	}

	return nil, err
}

func QueryToBytes(params *map[string]interface{}) (*bytes.Buffer, error) {
	bytesRepresentation, err := json.Marshal(*params)
	if err != nil {
		return nil, fmt.Errorf("can't transform bytes to query: %w", err)
	}
	return bytes.NewBuffer(bytesRepresentation), nil
}

func GetPayTypes(platformName string, c *config.Config) []string {
	var result []string

	for key, value := range c.PayTypes {
		if value {
			result = append(result, PayTypesDict[platformName][key])
		}
	}

	return result
}
