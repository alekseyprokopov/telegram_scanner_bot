package storage

type Configuration struct {
	UserName   string
	UserConfig Config
}

type Config struct {
	payTypes  []string
	minValue  int
	minSpread float64
	maxSpread float64
	platforms map[string]Platform
}

type Platform struct {
	IsActive     bool
	PlatformName string
	Roles        []Role
}

type Role struct {
	RoleName string
	IsActive bool
}

var DefaultConfig = Config{
	payTypes:  []string{"Sberbank", "Tinkoff", "QIWI", "YooMoney"},
	minValue:  10000,
	minSpread: 0.5,
	maxSpread: 10,
	platforms: map[string]Platform{
		"binance":
	},
}

func ()  {
	
}
