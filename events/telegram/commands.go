package telegram

import (
	"fmt"
	"log"
	"net/url"
	"scanner_bot/clients/telegramAPI"
	"scanner_bot/config"
	"scanner_bot/handler"
	"sort"
	"strconv"
	"strings"
)

const (
	SetConfigCmd = "/setConfig"

	ShowConfigCmd  = "Настройки"
	SetDefaultCmd  = "Сбросить настройки"
	BackCmd        = "Назад"
	SetLimitCmd    = "Лимит"
	SetSpreadCmd   = "Спред"
	SetOrdersCmd   = "Количество сделок"
	SetPaytypesCmd = "Способы оплаты"

	HelpCmd  = "/help"
	StartCmd = "/start"

	InsideTTCmd  = "Внутрибиржевые Т/Т"
	InsideTMCmd  = "Внутрибиржевые Т/М"
	OutsideTTCmd = "Межбиржевые Т/Т"
	OutsideTMCmd = "Межбиржевые Т/М"
	TestCmd      = "/test"
)

func (p *EventProcessor) doCmd(text string, chatID int64, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s, from user: %s", text, username)

	switch text {
	case SetDefaultCmd:
		return p.SetDefault(chatID)
	case HelpCmd:
		return p.SendHelp(chatID)
	case BackCmd:
		return p.MainMenu(chatID)
	case ShowConfigCmd:
		return p.ShowConfig(chatID)
	case StartCmd:
		return p.SaveConfig(chatID)

	//setting commands
	//case SetDefaultCmd:
	//	return p.SetDefault(chatID)

	case SetLimitCmd:
		return p.tg.SendInlineKeyboard(chatID, telegramAPI.LimitMessage, telegramAPI.LimitsKeyBoard)
	case SetSpreadCmd:
		return p.tg.SendInlineKeyboard(chatID, telegramAPI.SpreadMessage, telegramAPI.SpreadKeyBoard)

	case SetOrdersCmd:
		return p.tg.SendInlineKeyboard(chatID, telegramAPI.OrdersMessage, telegramAPI.OrdersKeyBoard)

	case SetPaytypesCmd:
		return p.tg.SendInlineKeyboard(chatID, telegramAPI.PaytypesMessage, telegramAPI.PaytypesKeyBoard)

		//search commands
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

func (p *EventProcessor) doCmdCallback(text string, chatID int64, username string) error {
	data := strings.Split(text, "_")
	typeData := data[0]
	textData := data[1]

	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}

	switch typeData {
	case "limit":
		conf.UserConfig.MinValue = textData
	case "order":
		conf.UserConfig.Orders, err = strconv.Atoi(textData)
	case "paytype":
		conf.UserConfig.PayTypes[textData] = !conf.UserConfig.PayTypes[textData]
	case "spread":
		conf.UserConfig.MinSpread, err = strconv.ParseFloat(textData, 64)
	}

	UserConfig, err := config.UserConfigToString(conf)
	if err != nil {
		return fmt.Errorf("can't convert userConfig to string")
	}

	p.storage.Update(chatID, UserConfig)
	p.ShowConfig(chatID)
	return nil
}

func (p *EventProcessor) SaveConfig(chatID int64) (err error) {
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
	if err := p.MainMenu(chatID); err != nil {
		return err
	}
	return nil
}

func (p *EventProcessor) SetDefault(chatID int64) error {
	conf := config.ToDefaultConfig(chatID)
	UserConfig, err := config.UserConfigToString(conf)
	if err != nil {
		return fmt.Errorf("can't convert userConfig to string")
	}
	if err := p.storage.Update(chatID, UserConfig); err != nil {
		return err
	}
	p.ShowConfig(chatID)
	return nil

}
func (p *EventProcessor) ShowConfig(chatID int64) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't show config (cmd): %w", err)
		}
	}()
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	result := msgConfig(conf)
	if err := p.tg.SendSettingsKeyboard(chatID, result); err != nil {
		return err
	}
	return nil

}
func (p *EventProcessor) MainMenu(chatID int64) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't back to main menu: %w", err)
		}
	}()
	if err := p.tg.SendMainKeyboard(chatID, "Главное меню"); err != nil {
		return err
	}
	return nil

}

func (p *EventProcessor) Test(chatID int64) error {
	conf, err := p.storage.GetConfig(chatID)
	if err != nil {
		return err
	}
	result, _ := p.handler.Platforms["huobi"].GetResult(conf)
	spot := result.Spot
	log.Printf("spot HUOBI: %+v", spot)
	return nil
}

func (p *EventProcessor) InsideTT(chatID int64) error {
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

	if len(data) < 10 {
		counter := 0
		msg := getResultMessage(data, &counter)
		if err := p.tg.SendMessage(chatID, msg); err != nil {
			return err
		}
		return nil
	}

	batches := getBatch(data)
	for _, batch := range batches {
		if err := p.tg.SendMessage(chatID, batch); err != nil {
			return err
		}

	}

	return nil
}

func (p *EventProcessor) InsideTM(chatID int64) error {
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

	if len(data) < 10 {
		counter := 0
		msg := getResultMessage(data, &counter)
		if err := p.tg.SendMessage(chatID, msg); err != nil {
			return err
		}
		return nil
	}

	batches := getBatch(data)
	for _, batch := range batches {
		if err := p.tg.SendMessage(chatID, batch); err != nil {
			return err
		}
	}
	return nil
}

func (p *EventProcessor) OutsideTT(chatID int64) error {
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

	if len(data) < 10 {
		counter := 0
		msg := getResultMessage(data, &counter)
		if err := p.tg.SendMessage(chatID, msg); err != nil {
			return err
		}
		return nil
	}
	batches := getBatch(data)
	for _, batch := range batches {
		if err := p.tg.SendMessage(chatID, batch); err != nil {
			return err
		}

	}

	return nil
}

func (p *EventProcessor) OutsideTM(chatID int64) error {
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

	if len(data) < 10 {
		counter := 0
		msg := getResultMessage(data, &counter)
		if err := p.tg.SendMessage(chatID, msg); err != nil {
			return err
		}
		return nil
	}

	batches := getBatch(data)
	for _, batch := range batches {
		if err := p.tg.SendMessage(chatID, batch); err != nil {
			return err
		}
	}
	return nil
}

func (p EventProcessor) getConfig() {

}

func (p *EventProcessor) SendHelp(chatId int64) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *EventProcessor) SendHello(chatId int64) error {
	return p.tg.SendMessage(chatId, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

func getBatch(data []handler.Chain) []string {
	counter := 1

	chunkSize := 10
	var result []string

	var first, last int
	for i := 0; i < len(data)/chunkSize+1; i++ {
		first = i * chunkSize
		last = i*chunkSize + chunkSize
		if last > len(data) {
			last = len(data)
		}
		if first == last {
			break
		}

		result = append(result, getResultMessage(data[first:last], &counter))
	}

	return result
}
