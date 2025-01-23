package repository

import (
	"context"
	"fmt"
	"log/slog"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
)

type Repository interface {
	GetRecommendations(ctx context.Context, userId int) ([]int, error)
	AddProductUpdate(ctx context.Context, product *myproto.ProductAction) error
	AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) error
}

// имплементация Repository интерфейса
type RecomRepository struct {
	relDB RelationalDataBase
	kvDB  KeyValueDatabse
	log   *slog.Logger
}

// слой репощитория - взаимодействие с Базами данных
func NewProductRepository(db *PostgresDB, kvDB *RedisRepository, log *slog.Logger) *RecomRepository {
	return &RecomRepository{
		relDB: db,
		kvDB:  kvDB,
		log:   log,
	}
}

func (r *RecomRepository) GetRecommendations(ctx context.Context, userId int) ([]int, error) {
	fi := "repository.RecomRepository.GetRecommendations"

	//проверяем есть ли в кэше продукты для пользователя
	prodictIds, err := r.kvDB.GetRecom(ctx, userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
	}
	//если есть возращаем
	if prodictIds != nil {
		r.log.Info(fmt.Sprintf("%s: Recom Get From redis UserID %d, Recommendations %v", fi, userId, prodictIds))
		return prodictIds, nil
	}

	//если в кэше нет, обращаемся в Базу
	prodictIds, err = r.relDB.GetProductsByUserId(ctx, userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}
	result := removeDuplicates(prodictIds)

	//добавляем в кэш полученную из реляцонной базы информацию
	err = r.kvDB.SetRecom(ctx, userId, result)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}
	r.log.Info("%s: Recom GetFrom Postgres", fi, fi)
	return result, nil
}

func (r *RecomRepository) AddProductUpdate(ctx context.Context, product *myproto.ProductAction) error {
	fi := "repository.RecomRepository.AddProductUpdate"

	err := r.relDB.AddProductUpdate(ctx, product)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}
	err = r.kvDB.DelAll(ctx)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func (r *RecomRepository) AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) error {
	fi := "repository.RecomRepository.AddUserUpdate"

	err := r.relDB.AddUserUpdate(ctx, user)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	//удаляем информацию о рекомендациях пользователя из кэша
	err = r.kvDB.DelRecom(ctx, int(user.UserId))
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}
	return nil
}

func removeDuplicates(slice []int) []int {
	keys := make(map[int]bool)
	var result []int

	for _, value := range slice {
		if _, exists := keys[value]; !exists {
			keys[value] = true
			result = append(result, value)
		}
	}

	return result
}
