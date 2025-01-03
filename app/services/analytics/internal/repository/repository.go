package repository

import (
	"log/slog"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
)

type Repository interface {
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

func (r *RecomRepository) AddProductUpdate(product *myproto.ProductAction) error {
	fi := "repository.RecomRepository.AddProductUpdate"

	timestamp, err := r.relDB.AddProductUpdate(product)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	if err := r.kwDB.SetProductUpdate(product, timestamp); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func (r *RecomRepository) AddUserUpdate(user *myproto.UserUpdate) error {
	fi := "repository.RecomRepository.AddUserUpdate"

	timestamp, err := r.relDB.AddUserUpdate(user)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	if err := r.kwDB.SetUserUpdate(user, timestamp); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}
