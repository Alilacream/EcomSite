package store

import "database/sql"

type CommandStore struct {
	db *sql.DB
}
