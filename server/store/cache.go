package store

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Set(ctx context.Context, key, set string) error
}
type CacheStore struct {
	db *redis.Client
}

func (s *CacheStore) Set(ctx context.Context, key string, quantity string) error {
	log.Println("Cache Hit!")
	return s.db.HSet(ctx, key, quantity).Err()
}

func (s *CacheStore) Get(ctx context.Context, key string, quantity string) (string, error) {
	name, err := s.db.HGet(ctx, key, quantity).Result()
	if err != nil {
		return "", err
	}
	return name, nil
}
