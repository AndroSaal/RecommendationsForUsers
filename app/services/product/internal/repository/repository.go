package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
)

// имплементация Repository интерфейса
type ProductRepository struct {
	relDB RelationalDataBase
	log   *slog.Logger
}

type Repository interface {
	AddNewProduct(ctx context.Context, productInfo *entities.ProductInfo) (int, error)
	UpdateProduct(ctx context.Context, productId int, productInfo *entities.ProductInfo) error
	DeleteProduct(ctx context.Context, productId int) error
}

// слой репощитория - взаимодействие с Базами данных
func NewProductRepository(db RelationalDataBase, log *slog.Logger) *ProductRepository {
	return &ProductRepository{
		relDB: db,
		log:   log,
	}
}

func (r *ProductRepository) AddNewProduct(ctx context.Context, productInfo *entities.ProductInfo) (int, error) {
	fi := "repository.ProductRepository.AddNewProduct"

	productId, err := r.relDB.AddNewProduct(ctx, productInfo)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return 0, err

	}
	r.log.Info(fmt.Sprintf("%s: add new product with id %d", fi, productId))
	return productId, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, productId int, productInfo *entities.ProductInfo) error {
	fi := "repository.ProductRepository.UpdateProduct"

	err := r.relDB.UpdateProduct(ctx, productId, productInfo)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}
	r.log.Info(fmt.Sprintf("%s: product with id %d updated", fi, productId))
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, productId int) error {
	fi := "repository.ProductRepository.DeleteProduct"

	err := r.relDB.DeleteProduct(ctx, productId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}
	r.log.Info(fmt.Sprintf("%s: product with id %d deleted", fi, productId))
	return nil
}
