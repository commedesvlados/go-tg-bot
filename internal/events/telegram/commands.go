package telegram

import (
	"context"
	"errors"
	"github.com/commedesvlados/go-tg-bot/internal/storage"
	"github.com/commedesvlados/go-tg-bot/pkg/lib/e"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	Helpcmd  = "/help"
	StartCmd = "/start"
)

func (p *EventProcessor) doCmd(ctx context.Context, text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new commad '%s' from '%s'\n", text, username)

	if isAddCmd(text) {
		return p.savePage(ctx, chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatId, username)
	case Helpcmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.sendUnknown(chatId)
	}
}

func (p *EventProcessor) savePage(ctx context.Context, chatId int, pageURL string, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return e.Wrap("page already exists: ", err)
	}
	if isExists {
		return p.tg.SendMessage(chatId, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return e.Wrap("can't save page", err)
	}

	if err := p.tg.SendMessage(chatId, msgSaved); err != nil {
		return e.Wrap("can't send message: ", err)
	}

	return nil
}

func (p *EventProcessor) sendRandom(ctx context.Context, chatId int, username string) error {
	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return e.Wrap("can't do command: can't send random", err)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return e.Wrap("can't do command: can't send random", err)
	}

	return p.storage.Remove(ctx, page)
}

func (p *EventProcessor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *EventProcessor) sendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}

func (p *EventProcessor) sendUnknown(chatId int) error {
	return p.tg.SendMessage(chatId, msgUnknownCommand)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
