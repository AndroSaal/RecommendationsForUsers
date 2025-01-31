package service

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/repository"
	"golang.org/x/net/context"
)

// имплементация интерфейса Service
type ProductService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewProductService(repo repository.Repository, log *slog.Logger) *ProductService {
	return &ProductService{
		repo: repo,
		log:  log,
	}
}

// функция вызывает метод репозитория по добавлению нового продукта
func (s *ProductService) CreateProduct(ctx context.Context, product *entities.ProductInfo) (int, error) {
	fi := "service.ProductService.CreateProduct"

	productId, err := s.repo.AddNewProduct(ctx, product)
	if err != nil {
		s.log.Error("%s: Error Creating Product: %v", fi, err)
		return 0, err
	}
	return productId, nil
}

// функция заменяет информацию о пользователе в базе по его id
func (s *ProductService) UpdateProduct(ctx context.Context, productId int, product *entities.ProductInfo) error {
	fi := "service.ProductService.UpdateProduct"

	if err := s.repo.UpdateProduct(ctx, productId, product); err != nil {
		s.log.Error("%s: Error Updating Product: %v", fi, err)
		return err

	}

	return nil
}

// функция удаляет информацию о существующем продукте
func (s *ProductService) DeleteProduct(ctx context.Context, productId int) error {
	fi := "service.ProductService.DeleteProduct"

	if err := s.repo.DeleteProduct(ctx, productId); err != nil {
		s.log.Error("%s: Error Deleting Product: %v", fi, err)
		return err
	}
	return nil
}
