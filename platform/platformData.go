package platform

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
			"Сбербанк": "185", "Тинькофф": "75", "Райффайзен": "64", "QIWI": "62", "ЮMoney": "274",
		},
		GarantexName: {
			"Сбербанк": "Сбербанк", "Тинькофф": "Тинькофф", "Райффайзен": "Райффайзен", "QIWI": "QIWI", "ЮMoney": "ЮMoney",
		},
	}
)

const (
	BinanceName   = "binance"
	P2PBinanceURL = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"

	ByBitName   = "bybit"
	P2PByBitURL = "https://api2.bybit.com/fiat/otc/item/online"

	HuobiName   = "huobi"
	P2PHuobiURL = "https://otc-akm.huobi.com/v1/data/trade-market?"

	GarantexName = "garantex"
)
