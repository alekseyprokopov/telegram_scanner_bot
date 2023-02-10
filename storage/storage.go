package storage

import (
	"encoding/json"
	"fmt"
)

type Configuration struct {
	UserId     int    `json:"user_name"`
	UserConfig Config `json:"user_config"`
}

type Config struct {
	PayTypes  []string            `json:"pay_types"`
	MinValue  int                 `json:"min_value"`
	MinSpread float64             `json:"min_spread"`
	MaxSpread float64             `json:"max_spread"`
	Platforms map[string]Platform `json:"platforms"`
}

type Platform struct {
	IsActivePlatform bool                    `json:"is_active_platform"`
	Roles            map[string]IsActiveRole `json:"roles"`
}

type IsActiveRole bool

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
	PayTypes:  []string{"Sberbank", "Tinkoff", "QIWI", "YooMoney"},
	MinValue:  10000,
	MinSpread: 0.5,
	MaxSpread: 10,
	Platforms: map[string]Platform{
		"binance": {
			IsActivePlatform: true,
			Roles: map[string]IsActiveRole{
				"taker/taker": true,
				"taker/maker": true,
				"maker/taker": true,
				"maker/maker": true,
			},
		},
	},
}

func toDefaultConfig(userId int) *Configuration {
	return &Configuration{
		UserId:     userId,
		UserConfig: *DefaultUserConfig,
	}
}
