package telegram

import (
	"fmt"
	"scanner_bot/config"
	"scanner_bot/handler"
	"strconv"
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
	userInfo := fmt.Sprintf("⚙Конфигурация пользователя: id%d⚙ \n\n", c.ChatId)

	minValue := fmt.Sprintf("Минимальное значение: %d \n", c.UserConfig.MinValue)
	minSpread := fmt.Sprintf("Минимальный спред: %.1f \n", c.UserConfig.MinSpread)
	maxSpread := fmt.Sprintf("Максимальный спред: %.1f \n", c.UserConfig.MaxSpread)
	payTypes := fmt.Sprintf("Банки: %s \n \n", payTypesToString(c))

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
		"🔴%s:\nПокупка: %s\nБанк: %s\nЦена: %f\nПродавец: %s\nЛимиты : %.1f - %.1f(%s)\nДоступно : %.2f(%s)\nСделки: %d\n",
		strings.ToUpper(buy.PlatformName),
		strings.ToUpper(buy.Asset),
		buy.BankName,
		buy.Cost,
		buy.SellerName,
		buy.MinLimit, buy.MaxLimit, buy.Fiat,
		buy.Available, buy.Asset,
		buy.SellerDeals,
	)

	spotInfo := fmt.Sprintf("\nПАРА: %s\nЦЕНА: %f\n\n", a.PairName, a.SpotPrice)
	if a.SpotName != "" {
		spotInfo = "\nСПОТ: " + a.SpotName + spotInfo
	}
	sellPlatformInfo := fmt.Sprintf(
		"🔴%s:\nПродажа: %s\nБанк: %s\nЦена: %f\nПродавец: %s\nЛимиты : %.1f - %.1f(%s)\nДоступно : %.2f(%s)\nСделки: %d\n",
		strings.ToUpper(sell.PlatformName),
		strings.ToUpper(sell.Asset),
		sell.BankName,
		sell.Cost,
		sell.SellerName,
		sell.MinLimit, sell.MaxLimit, sell.Fiat,
		sell.Available, sell.Asset,
		sell.SellerDeals,
	)

	profit := fmt.Sprintf("\nПРОФИТ: %.3f\n", a.Profit)

	result := buyPlatformInfo + spotInfo + sellPlatformInfo + profit
	return result
}

func getResultMessage(data []handler.Chain) string {
	resultMessage := ""
	for i, item := range data {
		chainMessage := "#" + strconv.Itoa(i) + "\n" + msgChain(&item)
		resultMessage += chainMessage
	}
	return resultMessage
}
