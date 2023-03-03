package events

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

//type Fetcher interface {
//	Fetch(limit int) ([]Event, error)
//}

type Fetcher interface {
	Fetch(limit int) (*tgbotapi.UpdatesChannel)
}

type Processor interface {
	Process(u tgbotapi.Update) error
}

type Type int

// Чтобы не было привязки к тг. Unknown - тип эвента, который мы не смогли распознать
const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
