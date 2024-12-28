package repository

import (
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/go-redis/redis"
)

type KeyValueDatabse interface {
	GetRecom(userId int) ([]int, error)
	SetRecom(userId int, productIds []int) error
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

func (r *RedisRepository) GetRecom(userId int) ([]int, error) {
	return nil, nil
}

func (r *RedisRepository) SetRecom(userId int, productIds []int) error {
	return nil
}
