package store

import "database/sql"

type ProductStore struct {
	db *sql.DB
}
