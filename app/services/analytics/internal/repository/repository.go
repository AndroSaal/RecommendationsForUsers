package repository

import (
	"context"
	"log/slog"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
)

type Repository interface {
	AddProductUpdate(ctx context.Context, product *myproto.ProductAction) error
	AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) error
}

// имплементация Repository интерфейса
type AnalyticsRepository struct {
	relDB RelationalDataBase
	kvdb  KeyValueDatabse
	log   *slog.Logger
}

// слой репощитория - взаимодействие с Базами данных
func NewAnalyticsRepository(db *PostgresDB, log *slog.Logger, kv *RedisRepository) *AnalyticsRepository {
	return &AnalyticsRepository{
		relDB: db,
		log:   log,
		kvdb:  kv,
	}
}

func (r *AnalyticsRepository) AddProductUpdate(ctx context.Context, product *myproto.ProductAction) error {
	fi := "repository.RecomRepository.AddProductUpdate"

	product.ProductKeyWords = removeDuplicates(product.ProductKeyWords)
	timestamp, err := r.relDB.AddProductUpdate(ctx, product)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	if err := r.kvdb.SetProductUpdate(ctx, product, timestamp); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func (r *AnalyticsRepository) AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) error {
	fi := "repository.RecomRepository.AddUserUpdate"

	user.UserInterests = removeDuplicates(user.UserInterests)
	timestamp, err := r.relDB.AddUserUpdate(ctx, user)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	if err := r.kvdb.SetUserUpdate(ctx, user, timestamp); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, value := range slice {
		if _, exists := keys[value]; !exists && value != "" {
			keys[value] = true
			result = append(result, value)
		}
	}

	return result
}
