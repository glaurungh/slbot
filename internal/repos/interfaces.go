package repos

import (
	"context"
	"github.com/glaurungh/slbot/internal/domain/models"
)

type StoreRepo interface {
	Put(context.Context, *models.Store) error
	GetByID(context.Context, int) (models.Store, error)
	GetByName(context.Context, string) (models.Store, error)
	GetAll(context.Context) ([]models.Store, error)
	Delete(context.Context, int) error
}

type ShoppingItemRepo interface {
	Put(context.Context, *models.ShoppingItem) error
	GetByStoreID(context.Context, int) ([]models.ShoppingItem, error)
	GetAll(context.Context) ([]models.ShoppingItem, error)
	Delete(context.Context, int) error
	DeleteMulti(context.Context, []int) error
}
