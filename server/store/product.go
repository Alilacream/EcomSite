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
	return nil, nil
}
