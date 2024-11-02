package repos

import (
	"context"
	"errors"
	"github.com/glaurungh/slbot/internal/domain/models"
	"sync"
)

// MockStoreRepo представляет моковую реализацию репозитория StoreRepo
type MockStoreRepo struct {
	data   map[int]models.Store
	mu     sync.RWMutex
	nextID int
}

// NewMockStoreRepo создает новый экземпляр MockStoreRepo
func NewMockStoreRepo() *MockStoreRepo {
	return &MockStoreRepo{
		data:   make(map[int]models.Store),
		nextID: 1,
	}
}

// Put добавляет или обновляет магазин в репозитории
func (repo *MockStoreRepo) Put(ctx context.Context, store *models.Store) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if store.ID == 0 {
		store.ID = repo.nextID
		repo.nextID++
	}

	repo.data[store.ID] = *store
	return nil
}

// GetByID возвращает магазин по его ID
func (repo *MockStoreRepo) GetByID(ctx context.Context, id int) (models.Store, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	store, exists := repo.data[id]
	if !exists {
		return models.Store{}, errors.New("store not found")
	}
	return store, nil
}

// GetByName возвращает магазин по его имени
func (repo *MockStoreRepo) GetByName(ctx context.Context, name string) (models.Store, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, store := range repo.data {
		if store.Name == name {
			return store, nil
		}
	}
	return models.Store{}, errors.New("store not found")
}

// GetAll возвращает все магазины
func (repo *MockStoreRepo) GetAll(ctx context.Context) ([]models.Store, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var stores []models.Store
	for _, store := range repo.data {
		stores = append(stores, store)
	}
	return stores, nil
}

// Delete удаляет магазин по его ID
func (repo *MockStoreRepo) Delete(ctx context.Context, id int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.data[id]; !exists {
		return errors.New("store not found")
	}
	delete(repo.data, id)
	return nil
}
