package messenger

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"

	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
	"github.com/lysenkopavlo/telegram-bot/internal/storage"
)

// commands are:
const (
	RndCmd       = "/rnd"
	HelpCmd      = "/help"
	StartCmd     = "/start"
	unableToSave = "can't save the page"
	unableToRnd  = "can't send the random page"
)

func (p *Processor) doCmd(chatID int, text, username string) error {
	text = strings.TrimSpace(text)

	slog.Info("Received command",
		slog.Group("User",
			slog.String("text is: ", text),
			slog.String("from user: ", username)))

	if isAddCommand(text) {
		slog.Info("Page is saved")
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		slog.Info("Random page has been sent")
		return p.sendRandom(chatID, username)
	case HelpCmd:
		slog.Info("Help info has been sent")
		return p.sendHelp(chatID)
	case StartCmd:
		slog.Info("Hello message has been sent")
		return p.sendHello(chatID)
	default:
		slog.Info("Received unknown command")
		return p.tgClient.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL, username string) error {

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	pageIsExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		slog.Error("IsExists func worked unexpectedly", "error", err)
		return e.WrapError(unableToSave, err)
	}
	if pageIsExists {
		slog.Info("Tried to save page existed page")
		return p.tgClient.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		slog.Error("Save func worked unexpectedly", "error: ", err)
		return e.WrapError(unableToSave, err)
	}

	if err := p.tgClient.SendMessage(chatID, msgSaved); err != nil {
		slog.Error("SendMessage func worked unexpectedly", "error: ", err)
		return e.WrapError(unableToSave, err)
	}

	slog.Info("savePage func worked as expected")
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) error {
	page, err := p.storage.PickRandom(context.Background(), username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {

		slog.Error("PickRandom func worked unexpectedly", "error", err)
		return e.WrapError(unableToRnd, err)
	}

	if errors.Is(err, storage.ErrNoSavedPages) {

		slog.Error("There is no save pages", "error", err)
		return p.tgClient.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tgClient.SendMessage(chatID, page.URL); err != nil {

		slog.Error("SendMessage func worked unexpectedly", "error", err)
		return e.WrapError(unableToRnd, err)
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tgClient.SendMessage(chatID, msgHelp)
}
func (p *Processor) sendHello(chatID int) error {
	return p.tgClient.SendMessage(chatID, msgHello)
}

func isAddCommand(text string) bool {
	return isItURL(text)
}

func isItURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
