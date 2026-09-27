package store

import (
	"context"
	"database/sql"

	"alilacream/ecom/internal/models"
)

type ProductRepository interface {
	GetbyID(ctx context.Context, id int64) (*models.Product, error)
}
type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) GetbyID(ctx context.Context, id int64) (*models.Product, error) {
	actualProduct := new(models.Product)
	query := `SELECT name, category , quantity FROM products WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&actualProduct.Category, &actualProduct.Name, &actualProduct.StockQuantity)
	if err != nil {
		return nil, err
	}
	return actualProduct, nil
}
