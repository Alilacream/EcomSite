package store

import (
	"context"
	"errors"
	"strconv"

	"alilacream/ecom/internal/models"

	"github.com/redis/go-redis/v9"
)

var ErrCartItemNotFound = errors.New("cart item not found")

// cart:{clientID} is a Redis Hash: field = ProductID, value = Quantity.
// A Hash (not a Set) is used because a cart item needs a quantity attached
// to each product, and HINCRBY gives atomic +/- updates without a
// read-modify-write race.
type CacheRepository interface {
	CreateItem(ctx context.Context, clientID string, item models.CartItem) (bool, error)
	GetItem(ctx context.Context, clientID string, productID int64) (*models.CartItem, error)
	GetCart(ctx context.Context, clientID string) ([]models.CartItem, error)
	UpdateItemQuantity(ctx context.Context, clientID string, productID int64, quantity int) error
	IncrementItemQuantity(ctx context.Context, clientID string, productID int64, delta int) (int64, error)
	RemoveItem(ctx context.Context, clientID string, productID int64) error
	ClearCart(ctx context.Context, clientID string) error
	CartExists(ctx context.Context, clientID string) (bool, error)
	ItemExists(ctx context.Context, clientID string, productID int64) (bool, error)
}

type CacheStore struct {
	db *redis.Client
}

func cartKey(clientID string) string {
	return "cart:" + clientID
}

func field(productID int64) string {
	return strconv.FormatInt(productID, 10)
}

// @CREATE
// Adds a new line to the cart. Uses HSETNX so it never clobbers an existing
// quantity — returns false if the product is already in the cart, in which
// case the caller should use UpdateItemQuantity or IncrementItemQuantity.
func (s *CacheStore) CreateItem(ctx context.Context, clientID string, item models.CartItem) (bool, error) {
	return s.db.HSetNX(ctx, cartKey(clientID), field(item.ProductID), item.Quantity).Result()
}

// @GET (single item)
func (s *CacheStore) GetItem(ctx context.Context, clientID string, productID int64) (*models.CartItem, error) {
	quantity, err := s.db.HGet(ctx, cartKey(clientID), field(productID)).Int()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCartItemNotFound
	}
	if err != nil {
		return nil, err
	}
	return &models.CartItem{ProductID: productID, Quantity: quantity}, nil
}

// @GET (whole cart)
func (s *CacheStore) GetCart(ctx context.Context, clientID string) ([]models.CartItem, error) {
	raw, err := s.db.HGetAll(ctx, cartKey(clientID)).Result()
	if err != nil {
		return nil, err
	}
	items := make([]models.CartItem, 0, len(raw))
	for productIDStr, quantityStr := range raw {
		productID, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, err
		}
		items = append(items, models.CartItem{ProductID: productID, Quantity: quantity})
	}
	return items, nil
}

// @UPDATE
// Overwrites the quantity of an existing line (e.g. a "set quantity to N"
// input from the UI). Fails with ErrCartItemNotFound if the product isn't
// already in the cart — use CreateItem for that.
func (s *CacheStore) UpdateItemQuantity(ctx context.Context, clientID string, productID int64, quantity int) error {
	exists, err := s.ItemExists(ctx, clientID, productID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCartItemNotFound
	}
	if quantity <= 0 {
		return s.RemoveItem(ctx, clientID, productID)
	}
	return s.db.HSet(ctx, cartKey(clientID), field(productID), quantity).Err()
}

// @UPDATE
// Atomically applies a relative change (+1 / -1, e.g. cart +/- buttons).
// A line that drops to zero or below is removed rather than left at 0.
func (s *CacheStore) IncrementItemQuantity(ctx context.Context, clientID string, productID int64, delta int) (int64, error) {
	key, f := cartKey(clientID), field(productID)
	newQuantity, err := s.db.HIncrBy(ctx, key, f, int64(delta)).Result()
	if err != nil {
		return 0, err
	}
	if newQuantity <= 0 {
		if err := s.db.HDel(ctx, key, f).Err(); err != nil {
			return 0, err
		}
		return 0, nil
	}
	return newQuantity, nil
}

// @DELETE (single item)
func (s *CacheStore) RemoveItem(ctx context.Context, clientID string, productID int64) error {
	return s.db.HDel(ctx, cartKey(clientID), field(productID)).Err()
}

// @DELETE (whole cart)
func (s *CacheStore) ClearCart(ctx context.Context, clientID string) error {
	return s.db.Del(ctx, cartKey(clientID)).Err()
}

func (s *CacheStore) CartExists(ctx context.Context, clientID string) (bool, error) {
	n, err := s.db.Exists(ctx, cartKey(clientID)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *CacheStore) ItemExists(ctx context.Context, clientID string, productID int64) (bool, error) {
	return s.db.HExists(ctx, cartKey(clientID), field(productID)).Result()
}
