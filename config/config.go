package config

import (
	"encoding/json"
	"fmt"
	"scanner_bot/platform"
)

type Configuration struct {
	ChatId     int    `json:"user_name"`
	UserConfig Config `json:"user_config"`
}

type Config struct {
	MinValue  int               `json:"min_value"`
	MinSpread float64           `json:"min_spread"`
	MaxSpread float64           `json:"max_spread"`
	Binance   platform.Platform `json:"binance"`
	Garantex  platform.Platform `json:"garantex"`
	Huobi     platform.Platform `json:"huobi"`
	ByBit     platform.Platform `json:"by_bit"`
}

func UserConfigToString(c *Configuration) (userConfig string, err error) {
	configuration := *c
	s, err := json.Marshal(configuration.UserConfig)
	if err != nil {
		return "", fmt.Errorf("can't encode config: %w", err)
	}
	return string(s), nil
}

func StringToConfig(userConfig string) (*Config, error) {
	pBytesUserConfig := []byte(userConfig)
	var data Config

	if err := json.Unmarshal(pBytesUserConfig, &data); err != nil {
		return nil, fmt.Errorf("can't transform string to Configuration: %w", err)
	}

	return &data, nil
}

var DefaultUserConfig = &Config{
	MinValue:  10000,
	MinSpread: 0.5,
	MaxSpread: 10,
	Binance: platform.Platform{
		PlatformName:     "binance",
		PayTypes:         []string{"RosBankNew", "TinkoffNew", "RaiffeisenBank", "QIWI", "YandexMoneyNew"},
		IsActivePlatform: true,
		Roles: map[string]bool{
			"taker/taker": true,
			"taker/maker": true,
			"maker/taker": false,
			"maker/maker": false},
	},
	Garantex: platform.Platform{
		PlatformName:     "garantex",
		PayTypes:         []string{"Cбер", "Тинькофф", "райффайзен", "QIWI", "ЮМани"},
		IsActivePlatform: true,
		Roles: map[string]bool{
			"taker/taker": true,
			"taker/maker": true,
			"maker/taker": false,
			"maker/maker": false,
		},
	},
	Huobi: platform.Platform{
		PlatformName:     "huobi",
		PayTypes:         []string{"29" /*SBER*/, "28" /*Tinkoff*/, "36" /*Райфайзен*/, "9" /*QIWI*/, "19" /*ЮМани*/},
		IsActivePlatform: true,
		Roles: map[string]bool{
			"taker/taker": true,
			"taker/maker": true,
			"maker/taker": false,
			"maker/maker": false,
		},
	},
	ByBit: platform.Platform{
		PlatformName:     "bybit",
		PayTypes:         []string{"185" /*SBER*/, "75" /*Tinkoff*/, "64" /*Райфайзен*/, "62" /*QIWI*/, "274" /*ЮМани*/},
		IsActivePlatform: true,
		Roles: map[string]bool{
			"taker/taker": true,
			"taker/maker": true,
			"maker/taker": false,
			"maker/maker": false,
		},
	},
}

func ToDefaultConfig(userId int) *Configuration {
	return &Configuration{
		ChatId:     userId,
		UserConfig: *DefaultUserConfig,
	}
}
