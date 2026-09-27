package store

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Product  ProductRepository
	Order    OrderRepository
	Customer CustomerRepository
	Cache    CacheRepository
}

func NewStorage(db *sql.DB, rclient *redis.Client) Storage {
	return Storage{
		Product:  &ProductStore{db},
		Order:    &OrderStore{db},
		Customer: &CustomerStore{db},
		Cache:    &CacheStore{rclient},
	}
}
