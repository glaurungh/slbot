package main

import (
	"context"
	_ "fmt"
	"github.com/glaurungh/slbot/internal/repos"
	"github.com/glaurungh/slbot/internal/services"
	"github.com/glaurungh/slbot/pkg/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func main() {
	// Bot creating
	telegramBot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	telegramBot.Debug = true

	// Подключение к PostgreSQL с использованием pgxpool
	dbURL := os.Getenv("DATABASE_URL") // ожидаем, что строка подключения будет в формате: postgres://user:password@host:port/dbname
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Unable to parse database URL:", err)
	}
	// Установка максимального количества соединений в пуле
	config.MaxConns = 10
	// Настройка времени ожидания для подключения
	config.MaxConnLifetime = 30 * time.Second
	config.HealthCheckPeriod = 2 * time.Second

	// Создание пула соединений
	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to establish database connection:", err)
	}
	defer dbPool.Close()

	// Создание репозитория и сервиса
	storeRepo := repos.NewPostgresStoreRepo(dbPool)
	itemRepo := repos.NewPostgresShoppingItemRepo(dbPool)
	storeService := services.NewStoreService(storeRepo)
	itemService := services.NewShoppingItemService(itemRepo)

	shoppingListBot := bot.NewBot(telegramBot, storeService, itemService)

	if err = shoppingListBot.Start(); err != nil {
		log.Fatal(err)
	}
}
