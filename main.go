package main

import (
	"flag"
	"log"
	tgClient "scanner_bot/clients/telegram"
	eventConsumer "scanner_bot/consumer/event-consumer"
	eventProcessor "scanner_bot/events/telegram"
	"scanner_bot/storage/sqlite"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "data/storage.db"
	batchSize   = 10000
)

func main() {
	s, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := eventProcessor.New(
		tgClient.New(tgBotHost, token()),
		s,
	)

	log.Print("service started...")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}

func token() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram",
	)
	//помещает значение в Токен
	flag.Parse()

	if *token == "" {
		log.Fatal("token is missing")
	}

	return *token
}
