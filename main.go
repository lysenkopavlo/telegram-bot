package main

import (
	"flag"
	"log"

	"github.com/lysenkopavlo/telegram-bot/internal/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {

	tgClient := telegram.New(tgBotHost, mustToken())
}

func mustToken() string {
	token := flag.String(
		"telegram-bot-token",
		"",
		"token to access telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}

	return *token
}
