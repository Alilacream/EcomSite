package db

import (
	"alilacream/ecom/config"

	"github.com/redis/go-redis/v9"
)

func RedisNew(db *config.DBConfig) (*redis.Client, error) {
	opts, err := redis.ParseURL(db.RedisURL)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opts)
	return rdb, nil
}
