package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/go-redis/redis"
)

type RedisRepository struct {
	KVDB *redis.Client
}

func NewRedisDB(cfg *config.KeyValueConfig) *RedisRepository {
	return &RedisRepository{
		KVDB: redis.NewClient(&redis.Options{
			Addr: cfg.Addr,
		}),
	}
}

func (r *RedisRepository) GetRecom(ctx context.Context, userId int) ([]int, error) {
	userIdKey := strconv.Itoa(userId)
	jsonData, err := r.KVDB.Get(userIdKey).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error getting product ids for user %d: %w", userId, err)
	}

	var productIds []int
	err = json.Unmarshal(jsonData, &productIds)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling product ids: %w", err)
	}
	return productIds, nil
}

func (r *RedisRepository) SetRecom(ctx context.Context, userId int, productIds []int) error {
	userIdKey := strconv.Itoa(userId)

	jsonData, err := json.Marshal(productIds)
	if err != nil {
		return fmt.Errorf("error marshalling product ids: %w", err)
	}

	err = r.KVDB.Set(userIdKey, jsonData, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error setting product ids for user %d: %w", userId, err)
	}

	return nil
}

func (r *RedisRepository) DelRecom(ctx context.Context, userId int) error {
	userIdKey := strconv.Itoa(userId)

	_, err := r.KVDB.Del(userIdKey).Result()

	if err != nil {
		return fmt.Errorf("error deleting product ids for user %d: %w", userId, err)
	}
	return nil
}

func (r *RedisRepository) DelAll(ctx context.Context) error {
	_, err := r.KVDB.FlushAll().Result()

	if err != nil {
		return fmt.Errorf("error deleting all product ids: %w", err)
	}
	return nil
}
