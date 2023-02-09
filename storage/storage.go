package storage

type Configuration struct {
	UserName   string `json:"user_name"`
	UserConfig Config `json:"user_config"`
}

type Config struct {
	payTypes  []string `json:"pay_types"`
	minValue  int `json:"min_value"`
	minSpread float64 `json:"min_spread"`
	maxSpread float64 `json:"max_spread"`
	platforms map[string]Platform `json:"platforms"`
}

type Platform struct {
	IsActivePlatform bool `json:"is_active_platform"`
	Roles    map[string]IsActiveRole `json:"roles"`
}

type IsActiveRole bool

var DefaultConfig = Config{
	payTypes:  []string{"Sberbank", "Tinkoff", "QIWI", "YooMoney"},
	minValue:  10000,
	minSpread: 0.5,
	maxSpread: 10,
	platforms: map[string]Platform{
		"binance": {
			IsActivePlatform: true,
			Roles: map[string]IsActiveRole{
				"taker/taker": true,
				"taker/maker":true,
				"maker/taker":true,
				"maker/maker":true,
			},
		},
	},
}
