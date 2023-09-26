package telegram

import (
	"bot-saver/clients/telegram"
	"bot-saver/storage"
)

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
