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

	if err := r.relDB.AddProductUpdate(product); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	if err := r.kwDB.SetProductUpdate(product); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func (r *RecomRepository) AddUserUpdate(user *myproto.UserUpdate) error {
	fi := "repository.RecomRepository.AddUserUpdate"

	if err := r.relDB.AddUserUpdate(user); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	if err := r.kwDB.SetUserUpdate(user); err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}
