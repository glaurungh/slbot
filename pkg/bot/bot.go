package bot

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
	// Устанавливаем offset = -1 для пропуска старых сообщений
	u := tgbotapi.NewUpdate(-1)
	u.Timeout = 60

	updatesChanel, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updatesChanel, nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery != nil {
			if err := b.handleCallbackQuery(update.CallbackQuery); err != nil {
				log.Println(err)
			}
			continue
		}
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

func (b *Bot) userState(userId int) string {
	return b.userStates[int64(userId)]
}
