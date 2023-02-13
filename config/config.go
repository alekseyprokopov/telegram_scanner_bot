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
		PlatformName:     "binace",
		PayTypes:         []string{"RosBankNew", "TinkoffNew", "QIWI", "YandexMoneyNew"},
		IsActivePlatform: true,
		Roles: map[string]bool{
			"taker/taker": true,
			"taker/maker": true,
			"maker/taker": false,
			"maker/maker": false},
	},
	Garantex: platform.Platform{
		PlatformName:     "garantex",
		PayTypes:         []string{"Cбер", "Тинькофф"},
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
