.PHONY:
.SILENT:

# Переменные окружения
DB_URL := $(shell echo $$DATABASE_URL) # URL базы данных из переменной окружения

# Путь к папке с миграциями
MIGRATIONS_DIR := migrations

# Цель по умолчанию
.DEFAULT_GOAL := help

build: ## Компиляция
	go build -o ./.bin/bot cmd/bot/main.go

run: build  ## Запуск
	./.bin/bot

help:  ## Отображает список доступных команд
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

migrate-new: ## Создать новую миграцию
ifndef name
	$(error "Укажите имя миграции с параметром name=<name>")
endif
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-all: ## Применить все миграции
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down-last: ## Откатить последнюю миграцию
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

migrate-force-set: ## Принудительно установить версию миграции (например, make force version=1)
ifndef version
	$(error "Укажите номер версии с параметром version=<номер версии>")
endif
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) force $(version)

migrate-reset: ## Откатить все миграции и применить заново
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up