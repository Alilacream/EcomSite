package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Set(ctx context.Context, key, set string) (bool, error)
}
type CacheStore struct {
	db *redis.Client
}

func (s *CacheStore) Set(ctx context.Context, key, set string) (bool, error) {
	return false, nil
}
