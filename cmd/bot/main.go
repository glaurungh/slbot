package main

import (
	_ "fmt"
	"github.com/glaurungh/slbot/pkg/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

func main() {
	telegramBot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	telegramBot.Debug = true

	shoppingListBot := bot.NewBot(telegramBot)

	if err = shoppingListBot.Start(); err != nil {
		log.Fatal(err)
	}
}
