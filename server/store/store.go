package store

import (
	"database/sql"
)

type Storage struct {
	Product  ProductRepository
	Order    OrderRepository
	Customer CustomerRepository
}

func NewPQStorage(db *sql.DB) Storage {
	return Storage{
		Product:  &ProductStore{db},
		Order:    &OrderStore{db},
		Customer: &CustomerStore{db},
	}
}
