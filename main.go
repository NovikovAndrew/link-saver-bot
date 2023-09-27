package main

import (
	"flag"
	"log"

	"bot-saver/clients/telegram"
	"bot-saver/storage/files"
)

func main() {
	token := mustFlag("telegram-token-bot", "", "for access telegram bot")
	host := mustFlag("telegram-host", "api.telegram.org", "for request to telegram bot")

	tgClient := telegram.New(host, token)
	storage := files.New("")
}

func mustFlag(name, value, usage string) string {
	flagStr := flag.String(name, value, usage)
	flag.Parse()
	if *flagStr == "" {
		log.Fatalf("flag is empty, name %s\n", name)
	}

	return *flagStr
}
