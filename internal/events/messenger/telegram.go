package messenger

import (
	"errors"
	"log/slog"

	"github.com/lysenkopavlo/telegram-bot/internal/clients/telegram"
	"github.com/lysenkopavlo/telegram-bot/internal/events"
	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
	"github.com/lysenkopavlo/telegram-bot/internal/storage"
)

type Processor struct {
	tgClient *telegram.Client
	offset   int
	storage  storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tgClient: client,
		//storage:  storage,
	}
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		slog.Info("Event is of type Message")
		return p.processMessage(event)
	default:
		slog.Error("unknown type of event")
		return e.WrapError("wrong type", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	slog.Info("Processing message...")

	var errMsg = "can't process a message"

	metaInfo, err := meta(event)
	if err != nil {
		slog.Error("meta func worked unexpectedly", "error: ", err)
		return e.WrapError(errMsg, err)
	}

	if err := p.doCmd(metaInfo.ChatID, event.Text, metaInfo.UserName); err != nil {
		slog.Error("doCmd func worked unexpectedly", "error: ", err)
		return e.WrapError(errMsg, err)
	}
	slog.Info("Message processed correctly")
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	slog.Info("Getting meta information of the message")
	if !ok {
		slog.Error(ErrUnknownMetaType.Error())
		return Meta{}, e.WrapError("can't get meta information", ErrUnknownMetaType)
	}

	slog.Info("Got meta information of the message")
	return res, nil
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tgClient.Updates(p.offset, limit)
	if err != nil {
		return nil, e.WrapError("Error while fetching: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		result = append(result, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return result, nil
}

func event(update telegram.Update) events.Event {
	updateType := fetchType(update)

	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	}

	return res
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
