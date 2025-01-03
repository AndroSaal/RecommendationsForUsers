package repository

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/entities"
	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/config"
	"github.com/go-redis/redis"
)

type KeyValueDatabse interface {
	SetUserUpdate(user *myproto.UserUpdate, timestamp time.Time) error
	SetProductUpdate(product *myproto.ProductAction, timestamp time.Time) error
	Del(id int) error
}

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

func (r *RedisRepository) SetProductUpdate(product *myproto.ProductAction, timestamp time.Time) error {
	productUpdate := entities.ProductFullUpdate{
		Product:   product,
		Timestamp: timestamp,
	}
	productIdKey := strconv.Itoa(int(productUpdate.Product.ProductId))

	jsonData, err := json.Marshal(productUpdate)
	if err != nil {
		return fmt.Errorf("error marshalling product ids: %w", err)
	}

	err = r.KVDB.Set(productIdKey, jsonData, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error setting product ids for user %s: %w", productIdKey, err)
	}

	return nil
}

func (r *RedisRepository) SetUserUpdate(user *myproto.UserUpdate, timestamp time.Time) error {
	userUpdate := entities.UserFullUpdate{
		User:      user,
		Timestamp: timestamp,
	}
	userIdKey := strconv.Itoa(int(userUpdate.User.UserId))

	jsonData, err := json.Marshal(userUpdate)
	if err != nil {
		return fmt.Errorf("error marshalling product ids: %w", err)
	}

	err = r.KVDB.Set(userIdKey, jsonData, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error setting product ids for user %s: %w", userIdKey, err)
	}

	return nil
}

func (r *RedisRepository) Del(id int) error {
	idKey := strconv.Itoa(id)

	_, err := r.KVDB.Del(idKey).Result()

	if err != nil {
		return fmt.Errorf("error deleting product ids for user %s: %w", idKey, err)
	}
	return nil
}
