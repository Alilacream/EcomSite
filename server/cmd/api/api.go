package main

import (
	"alilacream/ecom/config"
	"alilacream/ecom/store"
)

type application struct {
	config Config
	store  *store.Storage
}

type Config struct {
	addr string
	db   *config.DBConfig
}
