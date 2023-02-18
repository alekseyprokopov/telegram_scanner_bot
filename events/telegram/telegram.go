package telegram

import (
	"errors"
	"fmt"
	"scanner_bot/clients/telegram"
	"scanner_bot/events"
	"scanner_bot/handler"
	"scanner_bot/storage"
)

type EventProcessor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
	handler *handler.PlaftormHandler
}

type Meta struct {
	ChatId   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func New(client *telegram.Client, storage storage.Storage) *EventProcessor {
	return &EventProcessor{
		tg:      client,
		storage: storage,
		handler: handler.New(),
	}

}

func (p *EventProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}
	//заранее выделяем память для результата
	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, toEvent(u))
	}

	p.offset = updates[len(updates)-1].Id + 1

	return res, nil
}

func (p *EventProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can't process message: %w", ErrUnknownEventType)
	}

}

func (p *EventProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message %w", err)
	}
	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return fmt.Errorf("can't process message %w", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	//type assertion
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, fmt.Errorf("can't get meta: %w", ErrUnknownMetaType)
	}

	return res, nil
}

func toEvent(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			Username: upd.Message.From.Username,
			ChatId:   upd.Message.Chat.Id,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
