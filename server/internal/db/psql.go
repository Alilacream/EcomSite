package db

import (
	"alilacream/ecom/config"
	"context"
	"database/sql"
	"log"
	"time"
)

func PSQLNew(dbConf *config.DBConfig) (*sql.DB, error) {
	log.Printf("Connecting to: %s", dbConf.DSN)
	db, err := sql.Open("postgres", dbConf.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(dbConf.MaxIdleConnections)
	db.SetMaxOpenConns(dbConf.MaxOpenConnections)
	duration, err := time.ParseDuration(dbConf.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
