package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"scanner_bot/storage"
	"scanner_bot/storage/files"
	"strings"
)

const (
	SetConfigCmd  = "/setConfig"
	ShowConfigCmd = "/showConfig"
	HelpCmd       = "/help"
	StartCmd      = "/start"
)

func (p *EventProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s, from user: %s", text, username)

	if isAddCmd(text) {
		return p.SavePage(chatID, text, username)
	}

	switch text {
	case SetConfigCmd:
		return p.SetConfig(chatID, username)
	case HelpCmd:
		return p.SendHelp(chatID)

	case ShowConfigCmd:
		return p.ShowConfig(chatID)

	case StartCmd:
		return p.SendHello(chatID)

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *EventProcessor) SavePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't save page (cmd): %w", err)
		}
	}()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)

	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *EventProcessor) SetConfig(id int, username string) error {

}

func (p *EventProcessor) ShowConfig(id int, username string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't send random page (cmd): %w", err)
		}
	}()

	config, err := p.storage.Pick()

}

func (p *EventProcessor) sendRandom(chatId int, username string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't send random page (cmd): %w", err)
		}
	}()

	page, err := p.storage.Pick(username)

	if err != nil && !errors.Is(err, files.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, files.ErrNoSavedPages) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
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
