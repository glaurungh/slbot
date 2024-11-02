package repos

import (
	"context"
	"errors"
	"github.com/glaurungh/slbot/internal/domain/models"
	"sync"
)

type MockShoppingItemRepo struct {
	mu     sync.Mutex
	items  map[int]models.ShoppingItem
	lastID int
}

// Создание нового мокового репозитория
func NewMockShoppingItemRepo() *MockShoppingItemRepo {
	return &MockShoppingItemRepo{
		items: make(map[int]models.ShoppingItem),
	}
}

// Добавление или обновление товара
func (m *MockShoppingItemRepo) Put(ctx context.Context, item *models.ShoppingItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Если ID не задан, создаем новый ID
	if item.ID == 0 {
		m.lastID++
		item.ID = m.lastID
	}

	m.items[item.ID] = *item
	return nil
}

// Получение товаров по ID магазина
func (m *MockShoppingItemRepo) GetByStoreID(ctx context.Context, storeID int) ([]models.ShoppingItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var items []models.ShoppingItem
	for _, item := range m.items {
		if item.StoreID == storeID {
			items = append(items, item)
		}
	}
	return items, nil
}

// Получение всех товаров
func (m *MockShoppingItemRepo) GetAll(ctx context.Context) ([]models.ShoppingItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var items []models.ShoppingItem
	for _, item := range m.items {
		items = append(items, item)
	}
	return items, nil
}

// Удаление товара по ID
func (m *MockShoppingItemRepo) Delete(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.items[id]; !exists {
		return errors.New("item not found")
	}

	delete(m.items, id)
	return nil
}

// Удаление нескольких товаров по списку ID
func (m *MockShoppingItemRepo) DeleteMulti(ctx context.Context, ids []int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, id := range ids {
		if _, exists := m.items[id]; exists {
			delete(m.items, id)
		}
	}
	return nil
}
