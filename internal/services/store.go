package services

import (
	"context"
	"errors"
	"github.com/glaurungh/slbot/internal/domain/models"
	"github.com/glaurungh/slbot/internal/repos"
)

type StoreService struct {
	repo repos.StoreRepo
}

func NewStoreService(repo repos.StoreRepo) *StoreService {
	return &StoreService{repo: repo}
}

// Create добавляет новый магазин
func (s *StoreService) Create(ctx context.Context, store *models.Store) (models.Store, error) {
	// Вставка магазина в репозиторий
	if err := s.repo.Put(ctx, store); err != nil {
		return models.Store{}, err
	}

	return *store, nil
}

// Update обновляет данные существующего магазина
func (s *StoreService) Update(ctx context.Context, store *models.Store) (models.Store, error) {
	// Обновление магазина в репозитории
	if err := s.repo.Put(ctx, store); err != nil {
		return models.Store{}, err
	}

	return *store, nil
}

// GetByID возвращает магазин по его ID
func (s *StoreService) GetByID(ctx context.Context, id int) (models.Store, error) {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Store{}, err
	}
	return store, nil
}

// GetByName возвращает магазин по его имени
func (s *StoreService) GetByName(ctx context.Context, name string) (models.Store, error) {
	store, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return models.Store{}, err
	}
	return store, nil
}

// GetAll возвращает все магазины
func (s *StoreService) GetAll(ctx context.Context) ([]models.Store, error) {
	stores, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return stores, nil
}

// Delete удаляет магазин по его ID
func (s *StoreService) Delete(ctx context.Context, id int) error {
	// Проверка существования магазина
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("store not found")
	}

	// Удаление магазина
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
