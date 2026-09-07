package store

import (
	"alilacream/ecom/internal/models"
	"context"
	"database/sql"
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
