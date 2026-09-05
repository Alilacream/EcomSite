package store

import "database/sql"

type ClientStore struct {
	db *sql.DB
}
