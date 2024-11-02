package bot

import (
	"context"
	"fmt"
	"github.com/glaurungh/slbot/internal/domain/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
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

	switch b.userState(userID) {
	case "waiting_for_store_name":
		handleAddStore(b, chatID, userID, message.Text)
	case "waiting_for_store_selection":
		// Игнорируем текстовые сообщения, пока пользователь выбирает магазин
	case "waiting_for_item_name":
		handleAddItem(b, chatID, userID, message.Text)
	default:
		_, err := b.bot.Send(
			tgbotapi.NewMessage(
				chatID,
				"Используйте команду, чтобы начать. Например: /add_store, /add_item или /view_list",
			),
		)
		return err
	}

	return nil
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

// Универсальная обработка CallbackQuery с учетом действия
func (b *Bot) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) error {

	// Парсим действие и данные из callback data
	parts := strings.Split(callbackQuery.Data, ":")
	if len(parts) != 2 {
		_, err := b.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQuery.ID, "Неверный формат данных"))
		return err
	}
	action := parts[0]
	data := parts[1]

	// Обработка в зависимости от типа действия
	switch action {
	case "select_store":
		handleSelectStore(b, callbackQuery, data)
	default:
		_, err := b.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQuery.ID, "Неизвестное действие"))
		return err
	}
	return nil
}

// Карты для хранения данных
var shoppingList = make(map[int]models.ShoppingItem)

// Переменные для генерации ID
var itemIDCounter int

// Функции для работы с магазинами
func handleAddStore(b *Bot, chatID int64, userID int, name string) {
	if name == "" {
		b.bot.Send(tgbotapi.NewMessage(chatID, "Название магазина не может быть пустым. Введите название снова:"))
		return
	}
	newStore := models.Store{Name: name}
	b.storeService.Create(context.Background(), &newStore)
	b.userStates[int64(userID)] = "" // Сброс состояния
	b.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Магазин '%s' добавлен.", name)))
}

// Функции для работы с элементами списка покупок
func handleAddItem(b *Bot, chatID int64, userID int, itemName string) {
	if itemName == "" {
		b.bot.Send(tgbotapi.NewMessage(chatID, "Название товара не может быть пустым. Введите название снова:"))
		return
	}
	b.userStates[int64(userID)] = "waiting_for_store_selection"

	// Функция для создания inline-кнопок с разбиением по строкам
	createStoresInlineKeyboard := func(maxCharsPerRow int) tgbotapi.InlineKeyboardMarkup {
		var rows [][]tgbotapi.InlineKeyboardButton
		var row []tgbotapi.InlineKeyboardButton
		currentRowLength := 0

		stores, _ := b.storeService.GetAll(context.Background())

		for _, store := range stores {
			callbackData := fmt.Sprintf("select_store:%d", store.ID)
			button := tgbotapi.NewInlineKeyboardButtonData(store.Name, callbackData)

			// Проверяем, если текущая строка + новая кнопка превышает maxCharsPerRow
			if currentRowLength+len(store.Name) > maxCharsPerRow && len(row) > 0 {
				rows = append(rows, row)
				row = nil
				currentRowLength = 0
			}

			row = append(row, button)
			currentRowLength += len(store.Name) + 1 // Учитываем пробел между кнопками
		}

		// Добавляем последнюю строку кнопок
		if len(row) > 0 {
			rows = append(rows, row)
		}

		return tgbotapi.NewInlineKeyboardMarkup(rows...)
	}

	// Создаем инлайн-кнопки для выбора магазина с форматом "select_store:<storeID>"
	keyboard := createStoresInlineKeyboard(50)

	// Отправляем сообщение с инлайн-кнопками
	msg := tgbotapi.NewMessage(chatID, "Выберите магазин для этого товара:")
	msg.ReplyMarkup = keyboard
	b.bot.Send(msg)

	// Сохраняем промежуточное значение
	shoppingList[itemIDCounter+1] = models.ShoppingItem{Name: itemName}
	itemIDCounter++
}

// Просмотр списка покупок
func handleViewList(b *Bot, chatID int64) {
	var result strings.Builder

	stores, _ := b.storeService.GetAll(context.Background())

	for _, item := range shoppingList {
		storeName := "Неизвестно"
		for _, store := range stores {
			if store.ID == item.StoreID {
				storeName = store.Name
				break
			}
		}
		result.WriteString(fmt.Sprintf("%d, Товар: %s, Магазин: %s\n", item.ID, item.Name, storeName))
	}

	if result.Len() == 0 {
		result.WriteString("Ваш список покупок пуст.")
	}

	b.bot.Send(tgbotapi.NewMessage(chatID, result.String()))

}

// Обработка выбора магазина через callback data
func handleSelectStore(b *Bot, callbackQuery *tgbotapi.CallbackQuery, data string) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	storeID, err := strconv.Atoi(data)
	if err != nil {
		b.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQuery.ID, "Неверный выбор магазина"))
		return
	}

	// Извлекаем название товара и ID магазина из временных данных
	itemName := shoppingList[itemIDCounter].Name

	// Добавляем товар в список покупок
	shoppingList[itemIDCounter] = models.ShoppingItem{Name: itemName, StoreID: storeID}

	// Сброс состояния и временных данных
	b.userStates[int64(userID)] = ""
	//delete(userTempData, int64(userID))

	stores, _ := b.storeService.GetAll(context.Background())

	// Подтверждаем добавление товара
	storeName := "Неизвестно"
	for _, store := range stores {
		if store.ID == storeID {
			storeName = store.Name
			break
		}
	}
	b.bot.Send(
		tgbotapi.NewMessage(
			chatID,
			fmt.Sprintf("Товар '%s' добавлен в список покупок для магазина '%s'.", itemName, storeName),
		),
	)

	// Ответ на callback, чтобы убрать "загрузка"
	b.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQuery.ID, "Магазин выбран"))
}
