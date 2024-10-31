package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

const (
	commandStart    = "start"
	commandAddStore = "add_store"
	commandAddItem  = "add_item"
	commandViewList = "view_list"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	chatID := message.Chat.ID
	userID := message.From.ID
	state := b.userStates[int64(userID)]

	switch state {
	case "waiting_for_store_name":
		handleAddStore(b, chatID, userID, message.Text)
	case "waiting_for_item_name":
		handleAddItem(b, chatID, userID, message.Text)
	default:
		b.bot.Send(tgbotapi.NewMessage(chatID, "Используйте команду, чтобы начать. Например: /add_store или /add_item"))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	//msg.ReplyToMessageID = message.MessageID

	_, err := b.bot.Send(msg)
	return err

}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	chatID := message.Chat.ID
	userID := message.From.ID

	switch message.Command() {
	case commandStart:
		msg := tgbotapi.NewMessage(chatID, "Привет! 👋 Как оно?")
		_, err := b.bot.Send(msg)
		return err

	case commandAddStore:
		b.userStates[int64(userID)] = "waiting_for_store_name"
		_, err := b.bot.Send(tgbotapi.NewMessage(chatID, "Введите название магазина:"))
		return err
	case commandAddItem:
		b.userStates[int64(userID)] = "waiting_for_item_name"
		_, err := b.bot.Send(tgbotapi.NewMessage(chatID, "Введите название товара:"))
		return err
	case commandViewList:
		handleViewList(b, chatID)
	default:
		msg := tgbotapi.NewMessage(chatID, "Такого не умею 🤷")
		_, err := b.bot.Send(msg)
		return err
	}
	return nil
}

// Карты для хранения данных
var stores = make(map[int]Store)
var shoppingList = make(map[int]ShoppingItem)

// Переменные для генерации ID
var storeIDCounter, itemIDCounter int

// Функции для работы с магазинами
func handleAddStore(b *Bot, chatID int64, userID int, name string) {
	if name == "" {
		b.bot.Send(tgbotapi.NewMessage(chatID, "Название магазина не может быть пустым. Введите название снова:"))
		return
	}
	storeIDCounter++
	stores[storeIDCounter] = Store{ID: storeIDCounter, Name: name}
	b.userStates[int64(userID)] = "" // Сброс состояния
	b.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Магазин '%s' добавлен с ID %d.", name, storeIDCounter)))
}

// Функции для работы с элементами списка покупок
func handleAddItem(b *Bot, chatID int64, userID int, itemName string) {
	if itemName == "" {
		b.bot.Send(tgbotapi.NewMessage(chatID, "Название товара не может быть пустым. Введите название снова:"))
		return
	}
	b.userStates[int64(userID)] = "waiting_for_item_quantity"
	b.bot.Send(tgbotapi.NewMessage(chatID, "Введите количество:"))

	// Сохраняем промежуточное значение
	shoppingList[itemIDCounter+1] = ShoppingItem{Name: itemName}
}

// Просмотр списка покупок
func handleViewList(b *Bot, chatID int64) {
	var result strings.Builder

	for _, item := range shoppingList {
		storeName := "Неизвестно"
		for _, store := range stores {
			if store.ID == item.StoreID {
				storeName = store.Name
				break
			}
		}
		result.WriteString(fmt.Sprintf("Товар: %s, Магазин: %s\n", item.Name, storeName))
	}

	if result.Len() == 0 {
		b.bot.Send(tgbotapi.NewMessage(chatID, "Ваш список покупок пуст."))
	} else {
		b.bot.Send(tgbotapi.NewMessage(chatID, result.String()))
	}
}
