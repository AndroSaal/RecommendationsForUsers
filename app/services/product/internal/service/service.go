package service

import "github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"

type Service interface {
	ProductCreater
	ProductUpdater
	ProductDeleter
}

type ProductCreater interface {
	CreateProduct(user *entities.ProductInfo) (int, error)
}

type ProductUpdater interface {
	UpdateProduct(userId int, user *entities.ProductInfo) error
}

type ProductDeleter interface {
	DeleteProduct(productId int) error
}
