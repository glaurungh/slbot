package repos

import (
	"context"
	"github.com/glaurungh/slbot/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresShoppingItemRepo struct {
	db *pgxpool.Pool
}

func NewPostgresShoppingItemRepo(db *pgxpool.Pool) *PostgresShoppingItemRepo {
	return &PostgresShoppingItemRepo{db: db}
}

// Добавление или обновление товара
func (r *PostgresShoppingItemRepo) Put(ctx context.Context, item *models.ShoppingItem) error {
	if item.ID == 0 {
		return r.create(ctx, item)
	}
	// Команда выполняет вставку или обновляет `name`, если `id` уже существует
	query := `
		INSERT INTO shopping_items (id, name, store_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, store_id = EXCLUDED.store_id;
	`
	_, err := r.db.Exec(ctx, query, item.ID, item.Name, item.StoreID)
	return err
}

// Добавление товара
func (r *PostgresShoppingItemRepo) create(ctx context.Context, item *models.ShoppingItem) error {
	query := `
		INSERT INTO shopping_items (name, store_id)
		VALUES ($1, $2);
	`
	_, err := r.db.Exec(ctx, query, item.Name, item.StoreID)
	return err
}

// Получение всех товаров по ID магазина
func (r *PostgresShoppingItemRepo) GetByStoreID(ctx context.Context, storeID int) ([]models.ShoppingItem, error) {
	query := `SELECT id, name, store_id FROM shopping_items WHERE store_id = $1`
	rows, err := r.db.Query(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ShoppingItem
	for rows.Next() {
		var item models.ShoppingItem
		if err := rows.Scan(&item.ID, &item.Name, &item.StoreID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// Получение всех товаров
func (r *PostgresShoppingItemRepo) GetAll(ctx context.Context) ([]models.ShoppingItem, error) {
	query := `SELECT id, name, store_id FROM shopping_items`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ShoppingItem
	for rows.Next() {
		var item models.ShoppingItem
		if err := rows.Scan(&item.ID, &item.Name, &item.StoreID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// Удаление товара по ID
func (r *PostgresShoppingItemRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM shopping_items WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// Удаление нескольких товаров по списку ID
func (r *PostgresShoppingItemRepo) DeleteMulti(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	query := `DELETE FROM shopping_items WHERE id = ANY($1)`
	_, err := r.db.Exec(ctx, query, ids)
	return err
}
