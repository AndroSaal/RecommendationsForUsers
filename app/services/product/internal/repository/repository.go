package repository

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
)

// имплементация Repository интерфейса
type ProductRepository struct {
	relDB RelationalDataBase
	log   *slog.Logger
}

type Repository interface {
	AddNewProduct(productInfo *entities.ProductInfo) (int, error)
	UpdateProduct(productId int, productInfo *entities.ProductInfo) error
	DeleteProduct(productId int) error
}

// слой репощитория - взаимодействие с Базами данных
func NewProductRepository(db *PostgresDB, log *slog.Logger) *ProductRepository {
	return &ProductRepository{
		relDB: db,
		log:   log,
	}
}

func (r *ProductRepository) AddNewProduct(productInfo *entities.ProductInfo) (int, error) {
	fi := "repository.UserRepository.AddNewUser"

	userId, err := r.relDB.AddNewProduct(productInfo)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return 0, err

	}

	return userId, nil
}

func (r *ProductRepository) UpdateProduct(productId int, productInfo *entities.ProductInfo) error {
	fi := "repository.UserRepository.UpdateUser"

	err := r.relDB.UpdateProduct(productId, productInfo)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(productId int) error {
	fi := "repository.UserRepository.DeleteUser"

	err := r.relDB.DeleteProduct(productId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}
	return nil
}
