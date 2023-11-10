package main

import (
	"context"
	"flag"
	"log"
	"log/slog"

	"github.com/lysenkopavlo/telegram-bot/internal/clients/telegram"
	"github.com/lysenkopavlo/telegram-bot/internal/consumer/eventconsumer"
	"github.com/lysenkopavlo/telegram-bot/internal/events/messenger"
	"github.com/lysenkopavlo/telegram-bot/internal/storage/sqlite"
)

const (
	tgBotHost   = "api.telegram.org"
	batchSize   = 100
	storagePath = "pages.db"
)

func main() {

	s, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatal("database doesn't work", err)
	}

	err = s.Init(context.TODO())
	if err != nil {
		log.Fatal("database isn't initialized", err)
	}

	eventFetcherProcessor := messenger.New(
		telegram.New(tgBotHost, mustToken()),
		s)

	slog.Info("Service is up")

	eventConsumer := eventconsumer.New(100,
		eventFetcherProcessor,
		eventFetcherProcessor,
	)
	if err := eventConsumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token to access telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Println("Token is not specified")
	}

	return *token
}
