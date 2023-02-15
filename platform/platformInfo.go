package platform

var (
	binanceTokens    = []string{"USTD", "BTC", "BUSD", "BNB", "ETH", "SHIB"}
	binanceTradeType = []string{"BUY", "SELL"}

	bybitTokens    = []string{"USDT", "BTC", "ETH", "USDC"}
	bybitTradeType = []string{"1", "0"} //side 1- купить. 0 - продать

	huobiTokens = []string{
		"2",  /*USDT*/
		"1",  /*BTC*/
		"62", /*USDD*/
		"4",  /*HT*/
		"22", /*TRX*/
		"3",  /*ETH*/
		"7",  /*XRP*/
		"8",  /*LTC*/
	}
	huobiTradeType = []string{"sell", "buy"} //sell-купить, buy - продать

	payTypesDict = map[string]map[string]string{
		BinanceName: {
			"Сбербанк": "RosBankNew", "Тинькофф": "TinkoffNew", "Райффайзен": "RaiffeisenBank", "QIWI": "QIWI", "ЮMoney": "YandexMoneyNew",
		},
		HuobiName: {
			"Сбербанк": "29", "Тинькофф": "28", "Райффайзен": "36", "QIWI": "9", "ЮMoney": "19",
		},
		ByBitName: {
			"Сбербанк": "185", "Тинькофф": "75", "Райффайзен": "64", "QIWI": "62", "ЮMoney": "274",
		},
		GarantexName: {
			"Сбербанк": "Сбербанк", "Тинькофф": "Тинькофф", "Райффайзен": "Райффайзен", "QIWI": "QIWI", "ЮMoney": "ЮMoney",
		},
	}
)

const (
	BinanceName    = "binance"
	p2pBinanceHost = "p2p.binance.com"
	p2pBinancePath = "bapi/c2c/v2/friendly/c2c/adv/search"
	binanceAPI     = ""

	ByBitName    = "bybit"
	p2pByBitHost = "api2.bybit.com"
	p2pByBitPath = "spot/api/otc/item/list"
	byBitAPI     = ""

	HuobiName    = "huobi"
	p2pHuobiHost = "otc-akm.huobi.com"
	p2pHuobiPath = "v1/data/trade-market?"
	huobiAPI     = ""

	GarantexName = "garantex"
)
