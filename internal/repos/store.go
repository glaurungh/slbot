package repos // PostgresStoreRepo представляет реализацию репозитория StoreRepo для PostgreSQL
import (
	"context"
	"errors"
	"github.com/glaurungh/slbot/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStoreRepo struct {
	db *pgxpool.Pool
}

// NewPostgresStoreRepo создает новый экземпляр PostgresStoreRepo
func NewPostgresStoreRepo(db *pgxpool.Pool) *PostgresStoreRepo {
	return &PostgresStoreRepo{db: db}
}

// Put добавляет или обновляет магазин в базе данных
func (repo *PostgresStoreRepo) Put(ctx context.Context, store *models.Store) error {
	if store.ID == 0 {
		return repo.create(ctx, store)
	}
	// Команда выполняет вставку или обновляет `name`, если `id` уже существует
	err := repo.db.QueryRow(ctx,
		`INSERT INTO stores (name)
		 VALUES ($1, $2)
		 ON CONFLICT (id) 
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		store.ID, store.Name,
	).Scan(&store.ID)

	if err != nil {
		return err
	}
	return nil
}

// Put добавляет или обновляет магазин в базе данных
func (repo *PostgresStoreRepo) create(ctx context.Context, store *models.Store) error {
	// Команда выполняет вставку или обновляет `name`, если `id` уже существует
	err := repo.db.QueryRow(ctx,
		`INSERT INTO stores (name)
		 VALUES ($1)
		 RETURNING id`,
		store.Name,
	).Scan(&store.ID)

	if err != nil {
		return err
	}
	return nil
}

// GetByID возвращает магазин по его ID
func (repo *PostgresStoreRepo) GetByID(ctx context.Context, id int) (models.Store, error) {
	var store models.Store
	err := repo.db.QueryRow(ctx,
		"SELECT id, name FROM stores WHERE id = $1", id,
	).Scan(&store.ID, &store.Name)

	if err != nil {
		return models.Store{}, errors.New("store not found")
	}
	return store, nil
}

// GetByName возвращает магазин по его имени
func (repo *PostgresStoreRepo) GetByName(ctx context.Context, name string) (models.Store, error) {
	var store models.Store
	err := repo.db.QueryRow(ctx,
		"SELECT id, name FROM stores WHERE name = $1", name,
	).Scan(&store.ID, &store.Name)

	if err != nil {
		return models.Store{}, errors.New("store not found")
	}
	return store, nil
}

// GetAll возвращает все магазины
func (repo *PostgresStoreRepo) GetAll(ctx context.Context) ([]models.Store, error) {
	rows, err := repo.db.Query(ctx, "SELECT id, name FROM stores ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []models.Store
	for rows.Next() {
		var store models.Store
		if err := rows.Scan(&store.ID, &store.Name); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stores, nil
}

// Delete удаляет магазин по его ID
func (repo *PostgresStoreRepo) Delete(ctx context.Context, id int) error {
	result, err := repo.db.Exec(ctx, "DELETE FROM stores WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("store not found")
	}
	return nil
}
