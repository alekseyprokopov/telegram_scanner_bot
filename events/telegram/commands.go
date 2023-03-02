package telegram

import (
	"fmt"
	"log"
	"net/url"
	"scanner_bot/config"
	"sort"
	"strings"
)

const (
	SetConfigCmd  = "/setConfig"
	ShowConfigCmd = "/showConfig"
	HelpCmd       = "/help"
	StartCmd      = "/start"
	InsideTTCmd   = "/inside_tt"
	InsideTMCmd   = "/inside_tm"
	OutsideTTCmd  = "/outside_tt"
	OutsideTMCmd  = "/outside_tm"
	TestCmd       = "/test"
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

	case InsideTTCmd:
		return p.InsideTT(chatID)
	case InsideTMCmd:
		return p.InsideTM(chatID)
	case OutsideTTCmd:
		return p.OutsideTT(chatID)
	case OutsideTMCmd:
		return p.OutsideTM(chatID)
	case TestCmd:
		return p.Test(chatID)

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

func (p *EventProcessor) Test(chatID int) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	result, _ := p.handler.Platforms["huobi"].GetResult(conf)
	spot := result.Spot
	log.Printf("spot HUOBI: %+v", spot)
	return nil
}

func (p *EventProcessor) InsideTT(chatID int) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	data := *p.handler.InsideTT(conf)
	sort.Slice(data, func(i, j int) bool { return data[i].Profit > data[j].Profit })

	if len(data) == 0 {
		p.tg.SendMessage(chatID, "Связки не найдены...")
		return nil
	}

	message := getResultMessage(data)
	p.tg.SendMessage(chatID, message)

	fmt.Println("ДЛИНА: ", len(data))

	return nil
}

func (p *EventProcessor) InsideTM(chatID int) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	data := *p.handler.InsideTM(conf)
	sort.Slice(data, func(i, j int) bool { return data[i].Profit > data[j].Profit })

	if len(data) == 0 {
		p.tg.SendMessage(chatID, "Связки не найдены...")
		return nil
	}

	message := getResultMessage(data)
	p.tg.SendMessage(chatID, message)

	fmt.Println("ДЛИНА: ", len(data))

	return nil
}

func (p *EventProcessor) OutsideTT(chatID int) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	data := *p.handler.OutsideTT(conf)
	sort.Slice(data, func(i, j int) bool { return data[i].Profit > data[j].Profit })

	if len(data) == 0 {
		p.tg.SendMessage(chatID, "Связки не найдены...")
		return nil
	}
	message := getResultMessage(data)
	p.tg.SendMessage(chatID, message)

	return nil
}

func (p *EventProcessor) OutsideTM(chatID int) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	data := *p.handler.OutsideTM(conf)
	sort.Slice(data, func(i, j int) bool { return data[i].Profit > data[j].Profit })

	if len(data) == 0 {
		p.tg.SendMessage(chatID, "Связки не найдены...")
		return nil
	}
	message := getResultMessage(data)
	fmt.Println(message)
	p.tg.SendMessage(chatID, message)

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
