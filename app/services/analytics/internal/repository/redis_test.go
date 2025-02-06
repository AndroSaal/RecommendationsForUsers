package repository

import (
	"context"
	"os"
	"testing"
	"time"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Тестирование RedisRepository как KV databse
func TestRedisRepository_SetProductUpdatet_Correct(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	product := myproto.ProductAction{
		ProductId:       1,
		Action:          "update",
		ProductKeyWords: []string{"test"},
	}
	time := time.Now()

	err1 := redisConn.SetProductUpdate(context.Background(), &product, time)
	assert.NoError(t, err1)
}

func TestRedisRepository_SetUserUpdate_Correct(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	user := myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"update"},
	}
	time := time.Now()

	err1 := redisConn.SetUserUpdate(context.Background(), &user, time)
	assert.NoError(t, err1)
}

func TestRedisRepository_Del_Correct(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	var cfg config.KeyValueConfig = config.KeyValueConfig{
		Addr: addr,
	}

	redisConn := NewRedisDB(&cfg)

	err1 := redisConn.Del(context.Background(), 1)
	assert.NoError(t, err1)
}
