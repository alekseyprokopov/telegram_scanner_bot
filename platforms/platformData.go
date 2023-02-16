package platforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scanner_bot/config"
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
	BinanceName    = "binance"
	P2PBinanceHost = "p2p.binance.com"
	P2PBinancePath = "bapi/c2c/v2/friendly/c2c/adv/search"
	BinanceAPI     = ""

	ByBitName    = "bybit"
	P2PByBitHost = "api2.bybit.com"
	P2PByBitPath = "spot/api/otc/item/list"
	ByBitAPI     = ""

	HuobiName    = "huobi"
	P2PHuobiHost = "otc-akm.huobi.com"
	P2PHuobiPath = "v1/data/trade-market?"
	HuobiAPI     = ""

	GarantexName = "garantex"
)


func GetQuery(platformName string, c *config.Config, token string, tradeType string) (result *bytes.Buffer, err error) {
	switch platformName {
	case BinanceName:
		return binanceGetQuery(c, token, tradeType)
	case HuobiName:
		return huobiGetQuery(c, token, tradeType)
	case ByBitName:
		return bybitGetQuery(c, token, tradeType)
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