package telegramAPI

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Client struct {
	tg     *tgbotapi.BotAPI
	config *tgbotapi.UpdateConfig
}

func New(token string) *Client {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("can't create BOT")
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return &Client{
		tg:     bot,
		config: &updateConfig,
	}
}

func (c *Client) Updates() *tgbotapi.UpdatesChannel {
	//updates, err := c.tg.GetUpdates(*c.config)
	updates := c.tg.GetUpdatesChan(*c.config)

	//updates := c.tg.GetUpdates(*c.config)
	//if err != nil {
	//	return nil, fmt.Errorf("can't get updates: %w", err)
	//}
	return &updates
}
func (c *Client) SendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"
	if _, err := c.tg.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)

	}
	return nil
}


