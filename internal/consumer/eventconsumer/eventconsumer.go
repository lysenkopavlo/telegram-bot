package eventconsumer

import (
	"log/slog"
	"time"

	"github.com/lysenkopavlo/telegram-bot/internal/events"
)

type Consumer struct {
	batchSize int // how many events handled in one time
	fetcher   events.Fetcher
	processor events.Processor
}

func New(batchSize int, fetcher events.Fetcher, processor events.Processor) *Consumer {
	return &Consumer{
		batchSize: batchSize,
		fetcher:   fetcher,
		processor: processor,
	}

}

func (c *Consumer) Start() error {
	slog.Info("Consumer has started")
	// infinity cycle to get events
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			slog.Error("Trouble consuming an event", "error", err)
			continue
		}
		// in case no events to be caught
		if len(gotEvents) == 0 {
			// wait for 1 second
			slog.Info("There is no events. Waiting 1 second")
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents...); err != nil {
			slog.Error(err.Error())
			continue
		}
	}

}

func (c *Consumer) handleEvents(events ...events.Event) error {
	for _, event := range events {
		slog.Info("Got new:", "event", event.Text)

		if err := c.processor.Process(event); err != nil {
			slog.Error("Something goes wrong with event processing", "error:", err)

			continue
		}
	}
	return nil
}
