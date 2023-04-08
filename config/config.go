package config

import (
	"encoding/json"
	"fmt"
)

type Configuration struct {
	ChatId     int64  `json:"user_name"`
	UserConfig Config `json:"user_config"`
}

type Config struct {
	MinValue  string  `json:"min_value"`
	Orders    int     `json:"orders"`
	MinSpread float64 `json:"min_spread"`
	//MaxSpread float64         `json:"max_spread"`
	PayTypes map[string]bool `json:"pay_types"`
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
	MinValue:  "",
	MinSpread: 0.2,
	//MaxSpread: 10,
	Orders: 10,
	PayTypes: map[string]bool{
		//"Сбербанк": true, "Тинькофф": true, "Райффайзен": true, "QIWI": true, "ЮMoney": true,
		"Payeer": true, "Advcash": true,
	},
}

func ToDefaultConfig(userId int64) *Configuration {
	return &Configuration{
		ChatId:     userId,
		UserConfig: *DefaultUserConfig,
	}
}
