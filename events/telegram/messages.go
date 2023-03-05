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
	userInfo := fmt.Sprintf("⚙<b>Пользователь:</b> id%d \n\n", c.ChatId)

	minValue := fmt.Sprintf("<i>Лимит:</i> %s \n", c.UserConfig.MinValue)
	orders := fmt.Sprintf("<i>Количество сделок:</i> %d \n", c.UserConfig.Orders)
	minSpread := fmt.Sprintf("<i>MIN спред:</i> %.1f \n", c.UserConfig.MinSpread)
	//maxSpread := fmt.Sprintf("MAX спред: %.1f \n", c.UserConfig.MaxSpread)
	payTypes := fmt.Sprintf("<i>Способ оплаты:</i> %s \n \n", payTypesToString(c))

	var result strings.Builder
	result.WriteString(userInfo)
	result.WriteString(minValue)
	result.WriteString(minSpread)
	result.WriteString(orders)
	//result.WriteString(maxSpread)
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
		"%s <b>%s:</b>\nПокупка: %s\nБанк: %s\nЦена: %g\nПродавец: %s\nЛимиты : %g - %g %s\nДоступно : %g %s\nСделки: %d\n",
		platformBall[buy.PlatformName],
		strings.Title(buy.PlatformName),
		strings.ToUpper(buy.Asset),
		buy.BankName,
		buy.Cost,
		buy.SellerName,
		buy.MinLimit, buy.MaxLimit, buy.Fiat,
		buy.Available, buy.Asset,
		buy.SellerDeals,
	)

	spotInfo := fmt.Sprintf("\nПара: %s\nЦена: %g\n\n", a.PairName, a.SpotPrice)
	if a.SpotName != "" {
		spotInfo = "\nСпот: " + a.SpotName + spotInfo
	}
	sellPlatformInfo := fmt.Sprintf(
		"%s <b>%s:</b>\nПродажа: %s\nБанк: %s\nЦена: %g\nПродавец: %s\nЛимиты : %g - %g %s\nДоступно : %g %s\nСделки: %d\n",
		platformBall[sell.PlatformName],
		strings.Title(sell.PlatformName),
		strings.ToUpper(sell.Asset),
		sell.BankName,
		sell.Cost,
		sell.SellerName,
		sell.MinLimit, sell.MaxLimit, sell.Fiat,
		sell.Available, sell.Asset,
		sell.SellerDeals,
	)

	profit := fmt.Sprintf("\n<b>Профит:</b> %.2f", a.Profit) + `%` + "\n"

	result := buyPlatformInfo + spotInfo + sellPlatformInfo + profit
	return result
}

func getResultMessage(data []handler.Chain, counter *int) string {
	var resultMessage []string

	sep := "--------\n"
	for _, item := range data {
		chainMessage := "#" + strconv.Itoa(*counter) + "\n" + msgChain(&item)
		resultMessage = append(resultMessage, chainMessage)
		*counter += 1

	}
	return strings.Join(resultMessage, sep)
}

var platformBall = map[string]string{
	"binance":  "🔴",
	"bybit":    "\U0001F7E0",
	"huobi":    "\U0001F7E1",
	"garantex": "\U0001F7E2",
}
