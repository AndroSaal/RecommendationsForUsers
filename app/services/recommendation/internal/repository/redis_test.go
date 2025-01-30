package repository

import (
	"context"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Тестирование RedisRepository как KV databse
func TestRedisRepository_SetGet_GetRecom_Correct(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	err1 := redisConn.SetRecom(context.Background(), 1, []int{1, 2, 3})
	assert.NoError(t, err1)

	recom, err2 := redisConn.GetRecom(context.Background(), 1)
	assert.NoError(t, err2)
	assert.Equal(t, []int{1, 2, 3}, recom)
}

func TestRedisRepository_GetRecom_CorrectButNothing(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	recom, err2 := redisConn.GetRecom(context.Background(), 2)
	assert.NoError(t, err2)
	assert.Nil(t, recom)
}

func TestRedisRepository_SetRecom_Incorrect(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)
	err1 := redisConn.SetRecom(context.Background(), 1, nil)
	assert.NoError(t, err1)
}

func TestRedisRepository_SetRecom_Incorrect2(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)
	err1 := redisConn.SetRecom(context.Background(), 1, []int{1, 2, 3})
	assert.NoError(t, err1)

	err2 := redisConn.SetRecom(context.Background(), 1, []int{1, 2, 3})
	assert.NoError(t, err2)
}

func TestRedisRepository_DelRecom_Correct(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	err1 := redisConn.SetRecom(context.Background(), 1, []int{1, 2, 3})
	assert.NoError(t, err1)

	recom, err2 := redisConn.GetRecom(context.Background(), 1)
	assert.NoError(t, err2)
	assert.Equal(t, []int{1, 2, 3}, recom)

	err3 := redisConn.DelRecom(context.Background(), 1)
	assert.NoError(t, err3)
	recom, err4 := redisConn.GetRecom(context.Background(), 1)
	assert.NoError(t, err4)
	assert.Nil(t, recom)
}

func TestRedisRepository_DelRecom_Incorrect(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	err3 := redisConn.DelRecom(context.Background(), 100)
	assert.NoError(t, err3)
}

func TestRedisRepository_DelAll_Correct(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	err1 := redisConn.SetRecom(context.Background(), 1, []int{1, 2, 3})
	assert.NoError(t, err1)

	err2 := redisConn.SetRecom(context.Background(), 2, []int{1, 2, 3})
	assert.NoError(t, err2)

	err3 := redisConn.DelAll(context.Background())

	assert.NoError(t, err3)

	recom, err4 := redisConn.GetRecom(context.Background(), 1)
	assert.NoError(t, err4)
	assert.Nil(t, recom)
	recom, err5 := redisConn.GetRecom(context.Background(), 2)
	assert.NoError(t, err5)
	assert.Nil(t, recom)

}
