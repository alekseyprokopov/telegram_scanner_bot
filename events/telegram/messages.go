package telegram

import (
	"fmt"
	"scanner_bot/config"
	"scanner_bot/handler"
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
	maxSpread := fmt.Sprintf("*Максимальный спред:* %.1f \n", c.UserConfig.MaxSpread)
	payTypes := fmt.Sprintf("*Банки:* %s \n \n", payTypesToString(c))

	var result strings.Builder
	result.WriteString(userInfo)
	result.WriteString(minValue)
	result.WriteString(minSpread)
	result.WriteString(maxSpread)
	result.WriteString(payTypes)

	return result.String()
}

func payTypesToString(c *config.Configuration) string {
	data := c.UserConfig.PayTypes
	var result []string
	for key, isActive := range data {
		if isActive {
			result = append(result, key)
		}
	}
	return strings.Join(result, ", ")
}

func msgChain(a *handler.Chain) string {
	buy := a.Buy
	sell := a.Sell

	buyPlatformInfo := fmt.Sprintf(
		"*🔴%s🔴:*\n*Тип сделки:* %s\n*Банк:* %s\n*Цена:* %.2f\n*Продавец:* %s\n*Лимиты (%s):* %.1f - %.1f\n*Доступно (%s):* %.2f\n*Сделки:* %d\n",
		strings.ToUpper(buy.PlatformName),
		buy.TradeType,
		buy.BankName,
		buy.Cost,
		buy.SellerName,
		buy.Fiat, buy.MinLimit, buy.MaxLimit,
		buy.Asset, buy.Available,
		buy.SellerDeals,
	)

	spotInfo := fmt.Sprintf("\n*ПАРА:* %s\n*СПОТ:* %.3f\n\n", a.PairName, a.SpotPrice)

	sellPlatformInfo := fmt.Sprintf(
		"*🔴%s🔴:*\n*Тип сделки:* %s\n*Банк:* %s\n*Цена:* %.2f\n*Продавец:* %s\n*Лимиты (%s):* %.1f - %.1f\n*Доступно (%s):* %.2f\n*Сделки:* %d\n",
		strings.ToUpper(sell.PlatformName),
		sell.TradeType,
		sell.BankName,
		sell.Cost,
		sell.SellerName,
		sell.Fiat, sell.MinLimit, sell.MaxLimit,
		sell.Asset, sell.Available,
		sell.SellerDeals,
	)

	profit := fmt.Sprintf("\n*ПРОФИТ:* %.3f\n", a.Profit)

	result := buyPlatformInfo + spotInfo + sellPlatformInfo + profit
	return result
}
