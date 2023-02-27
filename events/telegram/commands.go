package telegram

import (
	"fmt"
	"log"
	"net/url"
	"scanner_bot/config"
	"strings"
)

const (
	SetConfigCmd  = "/setConfig"
	ShowConfigCmd = "/showConfig"
	HelpCmd       = "/help"
	StartCmd      = "/start"
	GetCoursesCmd = "/getCourses"
)

func (p *EventProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s, from user: %s", text, username)

	switch text {
	//case SetConfigCmd:
	//	return p.SetConfig(chatID, username)
	case HelpCmd:
		return p.SendHelp(chatID)

	case ShowConfigCmd:
		return p.ShowConfig(chatID)

	case StartCmd:
		return p.SaveConfig(chatID)

	case GetCoursesCmd:
		return p.InsideTakerTaker(chatID)

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *EventProcessor) SaveConfig(chatID int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't save page (cmd): %w", err)
		}
	}()

	conf := config.ToDefaultConfig(chatID)
	isExists, err := p.storage.IsExists(chatID)

	if err != nil {
		return err
	}

	if isExists {

		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(conf); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

//func (p *EventProcessor) SetConfig(id int, username string) error {
//
//}

func (p *EventProcessor) ShowConfig(chatID int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't show config (cmd): %w", err)
		}
	}()
	//исправить!!!!!!
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}

	//result, _ := config.UserConfigToString(conf)
	result := msgConfig(conf)
	if err := p.tg.SendMessage(chatID, result); err != nil {
		return err
	}
	return nil

}

func (p *EventProcessor) InsideTakerTaker(chatID int) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	data := p.handler.InsideTakerTaker(conf)

	if len(*data) == 0 {
		p.tg.SendMessage(chatID, "Связки не найдены...")
		return nil
	}

	for _, item := range *data {
		chainMessage := msgChain(&item)
		p.tg.SendMessage(chatID, chainMessage)

	}

	return nil
}

func (p EventProcessor) getConfig() {

}

func (p *EventProcessor) SendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *EventProcessor) SendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
