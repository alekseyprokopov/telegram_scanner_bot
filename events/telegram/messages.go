package telegram

import (
	"fmt"
	"scanner_bot/config"
	"scanner_bot/platform"
	"strings"
)

const (
	msgHelp = `описание бота, бла...бла...Хос - лох.`

	msgHello = `Hello, motherfuckers`

	msgUnknownCommand = "Unknown command"

	msgNoSavedPages = "You have no saved pages"

	msgSaved = "Saved!"

	msgAlreadyExists = "You have already have this page in your list"
)

func msgConfig(c *config.Configuration) string {
	userInfo := fmt.Sprintf("⚙*Конфигурация пользователя:* id%d⚙ \n\n", c.ChatId)

	minValue := fmt.Sprintf("*Минимальное значение:* %d \n", c.UserConfig.MinValue)
	minSpread := fmt.Sprintf("*Минимальный спред:* %.1f \n", c.UserConfig.MinSpread)
	maxSpread := fmt.Sprintf("*Максимальный спред:* %.1f \n \n", c.UserConfig.MaxSpread)
	binanceInfo := platformParser(&c.UserConfig.Binance)
	garantexInfo := platformParser(&c.UserConfig.Garantex)

	var result strings.Builder
	result.WriteString(userInfo)
	result.WriteString(minValue)
	result.WriteString(minSpread)
	result.WriteString(maxSpread)

	result.WriteString(binanceInfo)
	result.WriteString(garantexInfo)

	return result.String()
}

func platformParser(p *platform.Platform) string {
	platformInfo := fmt.Sprintf("_%s INFO:_ \n", strings.ToUpper(p.PlatformName))
	platformPay := fmt.Sprintf("*Платежные системы:* %s \n", strings.Join(p.PayTypes, ", "))
	platformRoles := fmt.Sprintf("*Роли:* %s \n \n", rolesParser(&p.Roles))
	return platformInfo + platformPay + platformRoles
}

func rolesParser(r *map[string]bool) string{
	var rString strings.Builder
	for key, value := range *r{
		if value {
		rString.WriteString(key + ", ")
		}
	}

	return rString.String()
}
