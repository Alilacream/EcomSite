package store

import (
	"context"
	"database/sql"

	"alilacream/ecom/internal/models"
)

type OrderRepository interface {
	GetbyID(ctx context.Context, id int64) (*models.Order, error)
}
type OrderStore struct {
	db *sql.DB
}

func (s *OrderStore) GetbyID(ctx context.Context, id int64) (*models.Order, error) {
	return nil, nil
}
