package telegram

import (
	"bot-saver/clients/telegram"
	e "bot-saver/package/error"
	"bot-saver/storage"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	RandomCmd = "/random"
	Help      = "/help"
	StartCmd  = "/start"
)

func (p *Processor) doCmd(chatID int, text, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command `%s` from `%s`\n", text, username)

	if isURL(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RandomCmd:
		return p.sendRandom(chatID, username)
	case Help:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL, username string) (err error) {
	defer func() { err = e.Wrap("can't to save the page", err) }()
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
		Created:  time.Now(),
	}

	isExist, err := p.storage.IsExist(page)
	telegramMessage := sendMessage(chatID, p.tg)

	if err != nil {
		return err
	}

	if isExist {
		return telegramMessage(msgAlreadyExist)
	}

	if err = p.storage.Save(page); err != nil {
		return err
	}

	if err = telegramMessage(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't to send the random page", err) }()

	page, err := p.storage.PickRandom(username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return nil
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func sendMessage(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
