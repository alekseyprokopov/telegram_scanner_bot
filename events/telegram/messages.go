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

	msgSaved = "Saved!"

	msgAlreadyExists = "You have already have this page in your list"
)

func msgConfig(c *config.Configuration) string {
	userInfo := fmt.Sprintf("⚙*Конфигурация пользователя:* id%d⚙ \n\n", c.ChatId)

	minValue := fmt.Sprintf("*Минимальное значение:* %d \n", c.UserConfig.MinValue)
	minSpread := fmt.Sprintf("*Минимальный спред:* %.1f \n", c.UserConfig.MinSpread)
	maxSpread := fmt.Sprintf("*Максимальный спред:* %.1f \n \n", c.UserConfig.MaxSpread)
	binanceInfo := platformParser(&c.UserConfig.Binance)
	byBitInfo := platformParser(&c.UserConfig.ByBit)
	huobiInfo := platformParser(&c.UserConfig.Huobi)
	garantexInfo := platformParser(&c.UserConfig.Garantex)

	var result strings.Builder
	result.WriteString(userInfo)
	result.WriteString(minValue)
	result.WriteString(minSpread)
	result.WriteString(maxSpread)

	result.WriteString(binanceInfo)
	result.WriteString(byBitInfo)
	result.WriteString(huobiInfo)
	result.WriteString(garantexInfo)

	return result.String()
}

func platformParser(p *platform.Platform) string {
	platformInfo := fmt.Sprintf("_%s:_ \n", strings.ToUpper(p.PlatformName))
	platformPay := fmt.Sprintf("*Платежные системы:* %s \n", strings.Join(p.PayTypes, ", "))
	platformRoles := fmt.Sprintf("*Роли:* %s \n \n", rolesParser(&p.Roles))
	return platformInfo + platformPay + platformRoles
}

func rolesParser(r *map[string]bool) string {
	var rString strings.Builder
	for key, value := range *r {
		if value {
			rString.WriteString(key + ", ")
		}
	}
	return rString.String()
}

func msgAdvertise(a *platform.Advertise) string {
	platformInfo := fmt.Sprintf("*%s:*\n", a.PlatformName)
	typeInfo := fmt.Sprintf("*Тип сделки:* %s\n", a.TradeType)
	bankInfo := fmt.Sprintf("*Банк:* %s\n", a.BankName)
	priceInfo := fmt.Sprintf("*Цена:* %.2f\n", a.Cost)
	sellerInfo := fmt.Sprintf("*Продавец:* %s\n", a.SellerName)
	limitsInfo := fmt.Sprintf("*Лимиты (%s):* %.1f - %.1f\n", a.Fiat, a.MinLimit, a.MaxLimit)
	amountInfo := fmt.Sprintf("*Доступно (%s):* %.2f\n", a.Asset, a.Available)
	dealsInfo := fmt.Sprintf("*Сделки:* %d ", a.SellerDeals)

	var result strings.Builder
	result.WriteString(platformInfo)
	result.WriteString(typeInfo)
	result.WriteString(bankInfo)
	result.WriteString(priceInfo)

	result.WriteString(sellerInfo)
	result.WriteString(limitsInfo)
	result.WriteString(amountInfo)
	result.WriteString(dealsInfo)

	return result.String()
}
