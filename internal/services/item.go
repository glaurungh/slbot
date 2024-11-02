package services

import (
	"context"
	"errors"
	"github.com/glaurungh/slbot/internal/domain/models"
	"github.com/glaurungh/slbot/internal/repos"
)

type ShoppingItemService struct {
	repo repos.ShoppingItemRepo
}

// NewShoppingItemService создает новый экземпляр ShoppingItemService
func NewShoppingItemService(repo repos.ShoppingItemRepo) *ShoppingItemService {
	return &ShoppingItemService{repo: repo}
}

// Create добавляет новый элемент в список покупок
func (s *ShoppingItemService) Create(ctx context.Context, item *models.ShoppingItem) (models.ShoppingItem, error) {
	if item.Name == "" {
		return models.ShoppingItem{}, errors.New("item name cannot be empty")
	}
	if item.StoreID <= 0 {
		return models.ShoppingItem{}, errors.New("invalid store ID")
	}

	err := s.repo.Put(ctx, item)
	if err != nil {
		return models.ShoppingItem{}, err
	}
	return *item, nil
}

// Update обновляет существующий элемент списка покупок
func (s *ShoppingItemService) Update(ctx context.Context, item *models.ShoppingItem) (models.ShoppingItem, error) {
	if item.ID <= 0 {
		return models.ShoppingItem{}, errors.New("invalid item ID")
	}
	if item.Name == "" {
		return models.ShoppingItem{}, errors.New("item name cannot be empty")
	}

	err := s.repo.Put(ctx, item)
	if err != nil {
		return models.ShoppingItem{}, err
	}
	return *item, nil
}

// GetAll возвращает все элементы списка покупок
func (s *ShoppingItemService) GetAll(ctx context.Context) ([]models.ShoppingItem, error) {
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetByStoreID возвращает все элементы по идентификатору магазина
func (s *ShoppingItemService) GetByStoreID(ctx context.Context, storeID int) ([]models.ShoppingItem, error) {
	if storeID <= 0 {
		return nil, errors.New("invalid store ID")
	}

	items, err := s.repo.GetByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Delete удаляет элемент списка покупок по идентификатору
func (s *ShoppingItemService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid item ID")
	}

	return s.repo.Delete(ctx, id)
}

// DeleteMulti удаляет несколько элементов списка покупок по списку идентификаторов
func (s *ShoppingItemService) DeleteMulti(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return errors.New("no item IDs provided")
	}

	return s.repo.DeleteMulti(ctx, ids)
}
