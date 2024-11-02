package bot

import (
	"github.com/glaurungh/slbot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	storeService *services.StoreService
	itemService  *services.ShoppingItemService
	userStates   map[int64]string
}

func NewBot(bot *tgbotapi.BotAPI, storeService *services.StoreService, itemService *services.ShoppingItemService) *Bot {
	return &Bot{
		bot:          bot,
		storeService: storeService,
		itemService:  itemService,
		userStates:   make(map[int64]string),
	}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	if err := b.setUpCommands(); err != nil {
		return err
	}

	updatesChanel := b.initUpdatesChannel()

	b.handleUpdates(updatesChanel)

	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	// Устанавливаем offset = -1 для пропуска старых сообщений
	u := tgbotapi.NewUpdate(-1)
	u.Timeout = 60

	updatesChanel := b.bot.GetUpdatesChan(u)

	return updatesChanel
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

func (b *Bot) userState(userId int64) string {
	return b.userStates[userId]
}

func (b *Bot) setUpCommands() error {
	// Установка команд бота
	commands := []tgbotapi.BotCommand{
		{Command: commandStart, Description: "Начать работу с ботом"},
		{Command: commandAddStore, Description: "Добавить магазин"},
		{Command: commandAddItem, Description: "Добавить в список"},
		{Command: commandDeleteItems, Description: "Удалить из списка"},
		{Command: commandViewList, Description: "Просмотреть списка покупок"},
		//{Command: "delete_store", Description: "Удалить магазин"},
	}

	// Устанавливаем команды в Telegram
	_, err := b.bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Fatalf("Ошибка установки команд: %v", err)
	}
	return nil
}
