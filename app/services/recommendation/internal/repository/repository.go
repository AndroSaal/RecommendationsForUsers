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
	kwDB  KeyValueDatabse
	log   *slog.Logger
}

// слой репощитория - взаимодействие с Базами данных
func NewProductRepository(db *PostgresDB, log *slog.Logger) *RecomRepository {
	return &RecomRepository{
		relDB: db,
		log:   log,
	}
}

func (r *RecomRepository) GetRecommendations(userId int) ([]int, error) {
	fi := "repository.RecomRepository.GetRecommendations"

	//проверяем есть ли в кэше продукты для пользователя
	prodictIds, err := r.kwDB.GetRecom(userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}
	//если есть возращаем
	if prodictIds != nil {
		return prodictIds, nil
	}

	//если в кэше нет, обращаемся в Базу
	prodictIds, err = r.relDB.GetProductsByUserId(userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}

	return prodictIds, nil
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
	return nil
}
