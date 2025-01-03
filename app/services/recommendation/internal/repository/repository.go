package repository

import (
	"log/slog"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
)

type Repository interface {
	GetRecommendations(userId int) ([]int, error)
	AddProductUpdate(product *myproto.ProductAction) error
	AddUserUpdate(user *myproto.UserUpdate) error
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

func (r *RecomRepository) GetRecommendations(userId int) ([]int, error) {
	fi := "repository.RecomRepository.GetRecommendations"

	//проверяем есть ли в кэше продукты для пользователя
	prodictIds, err := r.kvDB.GetRecom(userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}
	//если есть возращаем
	if prodictIds != nil {
		r.log.Info("%s: Recom Get From redis UserID %d, Recommendations %w", fi, userId, prodictIds)
		return prodictIds, nil
	}

	//если в кэше нет, обращаемся в Базу
	prodictIds, err = r.relDB.GetProductsByUserId(userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}
	result := removeDuplicates(prodictIds)

	//добавляем в кэш полученную из реляцонной базы информацию
	err = r.kvDB.SetRecom(userId, result)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}
	r.log.Info("%s: Recom GetFrom Postgres", fi, fi)
	return result, nil
}

func (r *RecomRepository) AddProductUpdate(product *myproto.ProductAction) error {
	fi := "repository.RecomRepository.AddProductUpdate"

	err := r.relDB.AddProductUpdate(product)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func (r *RecomRepository) AddUserUpdate(user *myproto.UserUpdate) error {
	fi := "repository.RecomRepository.AddUserUpdate"

	err := r.relDB.AddUserUpdate(user)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	//удаляем информацию о рекомендациях пользователя из кэша
	err = r.kvDB.DelRecom(int(user.UserId))
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
