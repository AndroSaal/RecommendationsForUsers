package service

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/repository"
)

// имплементация интерфейса Service
type ProductService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewUserService(repo repository.Repository, log *slog.Logger) *ProductService {
	return &ProductService{
		repo: repo,
		log:  log,
	}
}

// функция вызывает метод репозитория по добавлению нового продукта
func (s *ProductService) CreateProduct(product *entities.ProductInfo) (int, error) {
	return s.repo.AddNewProduct(product)
}

// функция заменяет информацию о пользователе в базе по его id
func (s *ProductService) UpdateProduct(productId int, product *entities.ProductInfo) error {
	return s.repo.UpdateProduct(productId, product)
}

func (s *ProductService) DeleteProduct(productId int) error {
	return s.repo.DeleteProduct(productId)
}
