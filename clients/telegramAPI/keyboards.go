package telegramAPI

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	LimitMessage = "Лимит"
	SpreadMessage = "MIN спред"
	PaytypesMessage = "Способы оплаты"
	OrdersMessage = "Количество сделок"

)

var mainKeyBoard = tgbotapi.NewReplyKeyboard(

	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Внутрибиржевые Т/Т"),
		tgbotapi.NewKeyboardButton("Межбиржевые Т/Т"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Внутрибиржевые Т/М"),
		tgbotapi.NewKeyboardButton("Межбиржевые Т/М"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Настройки"),
	),
)

var settingsKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Лимит"),
		tgbotapi.NewKeyboardButton("Спред"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Количество сделок"),
		tgbotapi.NewKeyboardButton("Способы оплаты"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сбросить настройки"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var LimitsKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("500", "limit_500"),
		tgbotapi.NewInlineKeyboardButtonData("1 000", "limit_1000"),
		tgbotapi.NewInlineKeyboardButtonData("5 000", "limit_5000"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("10 000", "limit_10000"),
		tgbotapi.NewInlineKeyboardButtonData("30 000", "limit_30000"),
		tgbotapi.NewInlineKeyboardButtonData("50 000", "limit_50000"),
	),

)

var SpreadKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("0.1%", "spread_0.1"),
		tgbotapi.NewInlineKeyboardButtonData("0.2%", "spread_0.2"),
		tgbotapi.NewInlineKeyboardButtonData("0.4%", "spread_0.4"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("0.5%", "spread_0.5"),
		tgbotapi.NewInlineKeyboardButtonData("0.8%", "spread_0.8"),
		tgbotapi.NewInlineKeyboardButtonData("1.0%", "spread_0.9"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1.5%", "spread_1.5"),
		tgbotapi.NewInlineKeyboardButtonData("2.0%", "spread_2.0"),
		tgbotapi.NewInlineKeyboardButtonData("3.0%", "spread_3.0"),
	),

)

var PaytypesKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Сбербанк", "paytype_Сбербанк"),
		tgbotapi.NewInlineKeyboardButtonData("Тинькофф", "paytype_Тинькофф"),
		tgbotapi.NewInlineKeyboardButtonData("Райффайзен", "paytype_Райффайзен"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("QIWI", "paytype_QIWI"),
		tgbotapi.NewInlineKeyboardButtonData("ЮMoney", "paytype_ЮMoney"),
	),

)
var OrdersKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("10", "paytype_Сбербанк"),
		tgbotapi.NewInlineKeyboardButtonData("30", "paytype_Тинькофф"),
		tgbotapi.NewInlineKeyboardButtonData("50", "paytype_Райффайзен"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("70", "paytype_QIWI"),
		tgbotapi.NewInlineKeyboardButtonData("100", "paytype_ЮMoney"),
		tgbotapi.NewInlineKeyboardButtonData("150", "paytype_ЮMoney"),
	),

)