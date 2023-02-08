package event_consumer

import (
	"log"
	"scanner_bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		getEvents, err := c.fetcher.Fetch(c.batchSize)

		if err != nil {
			log.Printf("ERR consumer: %s", err.Error())

			continue
		}

		if len(getEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.HandleEvents(getEvents); err != nil {
			log.Printf("ERR HandleEvents: %s", err.Error())
			continue
		}

	}

}

func (c Consumer) HandleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't hadle event: %s", err.Error())
			continue
		}
	}
	return nil
}
