package repository

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/stretchr/testify/assert"
)

type MockRelationaldatabase struct{}

func (m *MockRelationaldatabase) AddNewProduct(ctx context.Context, product *entities.ProductInfo) (int, error) {
	if product.Category == "некорректно" {
		return 0, errors.New("ошибка")
	}
	return 1, nil
}
func (m *MockRelationaldatabase) UpdateProduct(ctx context.Context, productId int, uproduct *entities.ProductInfo) error {
	if uproduct.Category == "некорректно" {
		return errors.New("ошибка")
	}
	return nil
}
func (m *MockRelationaldatabase) DeleteProduct(ctx context.Context, productId int) error {
	if productId == 0 {
		return errors.New("ошибка")
	}
	return nil
}

func TestRepository_AddNewProduct_Correct(t *testing.T) {
	repository := NewProductRepository(
		&MockRelationaldatabase{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	productId, err := repository.AddNewProduct(context.Background(), &entities.ProductInfo{
		Category: "корректно",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, productId)
}

func TestRepository_AddNewProduct_CorrectButSomeError(t *testing.T) {
	repository := NewProductRepository(
		&MockRelationaldatabase{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	productId, err := repository.AddNewProduct(context.Background(), &entities.ProductInfo{
		Category: "некорректно",
	})

	assert.Error(t, err)
	assert.Equal(t, 0, productId)
}

func TestRepository_UpdateProduct_Correct(t *testing.T) {
	repository := NewProductRepository(
		&MockRelationaldatabase{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)
	err := repository.UpdateProduct(context.Background(), 1, &entities.ProductInfo{
		Category: "корректно",
	})

	assert.NoError(t, err)
}

func TestRepository_UpdateProduct_CorrectButSomeError(t *testing.T) {
	repository := NewProductRepository(
		&MockRelationaldatabase{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	err := repository.UpdateProduct(context.Background(), 1, &entities.ProductInfo{
		Category: "некорректно",
	})

	assert.Error(t, err)
}

func TestRepository_DeleteProduct_Correct(t *testing.T) {
	repository := NewProductRepository(
		&MockRelationaldatabase{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	err := repository.DeleteProduct(context.Background(), 1)

	assert.NoError(t, err)
}

func TestRepository_DeleteProductt_CorrectButSomeError(t *testing.T) {
	repository := NewProductRepository(
		&MockRelationaldatabase{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	err := repository.DeleteProduct(context.Background(), 0)

	assert.Error(t, err)
}
