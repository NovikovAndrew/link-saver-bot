package telegram

import (
	"bot-saver/clients/telegram"
	"bot-saver/events"
	e "bot-saver/package/error"
	"bot-saver/storage"
	"errors"
)

var (
	ErrUnknownType     = errors.New("can't provide unknown type")
	ErrUnknownMetaType = errors.New("can't to provide unknown meta type")
)

type Meta struct {
	ChatID   int
	Username string
}

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(tg *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tg,
		storage: storage,
	}
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknownType
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't to provide meta", err) }()
	meta, err := meta(event)

	if err != nil {
		return err
	}

	if err := p.doCmd(meta.ChatID, event.Text, meta.Username); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}

	return meta, nil
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)

	if err != nil {
		return nil, e.Wrap("can't to fetch updates", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(update telegram.Update) events.Event {
	eType := fetchType(update)
	event := events.Event{
		Type: eType,
		Text: fetchText(update),
	}

	if eType == events.Message {
		event.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return event
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
