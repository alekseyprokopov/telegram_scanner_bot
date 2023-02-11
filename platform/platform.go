package platform

type Platform struct {
	PlatformName     string
	PayTypes         []string        `json:"pay_types"`
	IsActivePlatform bool            `json:"is_active_platform"`
	Roles            map[string]bool `json:"roles"`
}

func New(platformName string, isActive bool, payTypes []string, roles map[string]bool) *Platform {
	return &Platform{
		PlatformName:     platformName,
		PayTypes:         payTypes,
		IsActivePlatform: isActive,
		Roles:            roles,
	}
}
