package store

import (
	"context"
	"database/sql"

	"alilacream/ecom/internal/models"
)

type CustomerRepository interface {
	GetbyID(ctx context.Context, id string) (*models.Customer, error)
}

type CustomerStore struct {
	db *sql.DB
}

func (s *CustomerStore) GetbyID(ctx context.Context, id string) (*models.Customer, error) {
	return nil, nil
}
