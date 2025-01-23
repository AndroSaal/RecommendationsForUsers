package service

import (
	"context"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
)

type Service interface {
	ProductCreater
	ProductUpdater
	ProductDeleter
}

type ProductCreater interface {
	CreateProduct(ctx context.Context, user *entities.ProductInfo) (int, error)
}

type ProductUpdater interface {
	UpdateProduct(ctx context.Context, userId int, user *entities.ProductInfo) error
}

type ProductDeleter interface {
	DeleteProduct(ctx context.Context, productId int) error
}
