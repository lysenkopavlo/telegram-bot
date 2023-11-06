package messenger

import "github.com/lysenkopavlo/telegram-bot/internal/clients/telegram"

type Processor struct {
	tgClient *telegram.Client
	offset   int
}
