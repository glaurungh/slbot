package main

import (
	_ "fmt"
	"github.com/glaurungh/slbot/pkg/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramBot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	telegramBot.Debug = true

	shoppingListBot := tgbot.NewBot(telegramBot)

	if err = shoppingListBot.Start(); err != nil {
		log.Fatal(err)
	}
}
