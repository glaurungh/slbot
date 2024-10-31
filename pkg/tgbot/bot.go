package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	userStates map[int64]string
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot:        bot,
		userStates: make(map[int64]string),
	}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updatesChanel, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updatesChanel)

	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updatesChanel, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updatesChanel, nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				log.Println(err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			log.Println(err)
		}

	}
}
