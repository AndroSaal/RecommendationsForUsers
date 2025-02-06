package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type RepositoryMock struct{}

func (m *RepositoryMock) AddNewProduct(ctx context.Context, productInfo *entities.ProductInfo) (int, error) {
	if productInfo.Category == "некорректно" {
		return 0, errors.New("ошибка")
	}
	return 1, nil
}

func (m *RepositoryMock) UpdateProduct(ctx context.Context, productId int, productInfo *entities.ProductInfo) error {
	if productInfo.Category == "некорректно" {
		return errors.New("ошибка")
	}
	return nil
}

func (m *RepositoryMock) DeleteProduct(ctx context.Context, productId int) error {
	if productId == 0 {
		return errors.New("ошибка")
	}
	return nil
}

func TestService_CreateProduct_Correct(t *testing.T) {
	service := NewProductService(
		&RepositoryMock{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	productId, err := service.CreateProduct(context.Background(), &entities.ProductInfo{
		Category: "корректно",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, productId)
}

func TestService_CreateProduct_CorrectButSomeError(t *testing.T) {
	service := NewProductService(
		&RepositoryMock{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	productId, err := service.CreateProduct(context.Background(), &entities.ProductInfo{
		Category: "некорректно",
	})

	assert.Error(t, err)
	assert.Equal(t, 0, productId)
}

func TestService_UpdateProduct_Correct(t *testing.T) {
	service := NewProductService(
		&RepositoryMock{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)
	err := service.UpdateProduct(context.Background(), 1, &entities.ProductInfo{
		Category: "корректно",
	})

	assert.NoError(t, err)
}

func TestService_UpdateProduct_CorrectButSomeError(t *testing.T) {
	service := NewProductService(
		&RepositoryMock{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	err := service.UpdateProduct(context.Background(), 1, &entities.ProductInfo{
		Category: "некорректно",
	})

	assert.Error(t, err)
}

func TestService_DeleteProduct_Correct(t *testing.T) {
	service := NewProductService(
		&RepositoryMock{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	err := service.DeleteProduct(context.Background(), 1)

	assert.NoError(t, err)
}

func TestService_DeleteProductt_CorrectButSomeError(t *testing.T) {
	service := NewProductService(
		&RepositoryMock{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	err := service.DeleteProduct(context.Background(), 0)

	assert.Error(t, err)
}
