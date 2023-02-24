package handler

var Binance = PlatformInfo{
	name:       "binance",
	p2pURL:     "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search",
	tokens:     []string{"USDT", "BTC", "BUSD", "BNB", "ETH", "SHIB"},
	tradeTypes: []string{"BUY", "SELL"},
	payTypesDict: map[string]string{
		"Сбербанк": "RosBankNew", "Тинькофф": "TinkoffNew", "Райффайзен": "RaiffeisenBank", "QIWI": "QIWI", "ЮMoney": "YandexMoneyNew",
		"RosBankNew": "Сбербанк", "TinkoffNew": "Тинькофф", "RaiffeisenBank": "Райффайзен", "YandexMoneyNew": "ЮMoney",
	},
}

var Bybit = PlatformInfo{
	name:       "bybit",
	p2pURL:     "https://api2.bybit.com/fiat/otc/item/online",
	tokens:     []string{"USDT", "BTC", "ETH", "USDC"},
	tradeTypes: []string{"1", "0"},
	payTypesDict: map[string]string{
		"185": "Сбербанк", "75": "Тинькофф", "64": "Райффайзен", "62": "QIWI", "274": "ЮMoney",
		"Сбербанк": "185", "Тинькофф": "75", "Райффайзен": "64", "QIWI": "62", "ЮMoney": "274",
	},
}

var Huobi = PlatformInfo{
	name:   "huobi",
	p2pURL: "https://otc-akm.huobi.com/v1/data/trade-market",
	tokens: []string{
		"2",  /*USDT*/
		"1",  /*BTC*/
		"62", /*USDD*/
		"4",  /*HT*/
		"22", /*TRX*/
		"3",  /*ETH*/
		"7",  /*XRP*/
		"8",  /*LTC*/
	},
	tradeTypes: []string{"BUY", "SELL"},
	tokensDict: map[string]string{
		"2":  "USDT",
		"1":  "BTC",
		"62": "USDD",
		"4":  "HT",
		"22": "TRX",
		"3":  "ETH",
		"7":  "XRP",
		"8":  "LTC",
	},
	payTypesDict: map[string]string{
		"Сбербанк": "29", "Тинькофф": "28", "Райффайзен": "36", "QIWI": "9", "ЮMoney": "19",
		"29": "Сбербанк", "28": "Тинькофф", "36": "Райффайзен", "9": "QIWI", "19": "ЮMoney",
	},
}
var allTokens = map[string]bool{
	"USDT": true,
	"BTC":  true,
	"BUSD": true,
	"BNB":  true,
	"ETH":  true,
	"SHIB": true,
	"USDC": true,
	"USDD": true,
	"HT":   true,
	"TRX":  true,
	"XRP":  true,
	"LTC":  true,
}
var allPairs = *CreatePairsSet(allTokens)

type PlatformInfo struct {
	name         string
	p2pURL       string
	tokens       []string
	tokensDict   map[string]string //необязательно
	tradeTypes   []string
	payTypesDict map[string]string
}

func CreatePairsSet(data map[string]bool) *map[string]bool {
	set := map[string]bool{}
	for key1, value1 := range data {
		for key2, value2 := range data {
			if value1 && value2 && key1 != key2 {
				set[key1+key2] = true
				set[key2+key1] = true
			}
		}
	}
	return &set
}
