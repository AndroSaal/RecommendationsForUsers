package repository

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
)

// имплементация Repository интерфейса
type UserRepository struct {
	relDB RelationalDataBase
	log   *slog.Logger
}

type Repository interface {
	AddNewProduct(productInfo *entities.ProductInfo) (int, error)
	UpdateProduct(productId int, productInfo *entities.ProductInfo) error
	DeleteProduct(productId int) error
}

// слой репощитория - взаимодействие с Базами данных
func NewUserRepository(db *PostgresDB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		relDB: db,
		log:   log,
	}
}

func (r *UserRepository) AddNewProduct(productInfo *entities.ProductInfo) (int, error) {
	fi := "repository.UserRepository.AddNewUser"

	userId, err := r.relDB.AddNewProduct(productInfo)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return 0, err

	}

	return userId, nil
}

func (r *UserRepository) UpdateProduct(productId int, productInfo *entities.ProductInfo) error {
	fi := "repository.UserRepository.UpdateUser"

	err := r.relDB.UpdateProduct(productId, productInfo)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}
