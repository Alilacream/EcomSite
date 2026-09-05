package store

import "database/sql"

type Storage struct {
	Product interface{}

	Command interface{}

	Client interface{}
}

func NewPQStorage(db *sql.DB) Storage {
	return Storage{
		Product: &ProductStore{db},
		Command: &CommandStore{db},
		Client:  &ClientStore{db},
	}
}
