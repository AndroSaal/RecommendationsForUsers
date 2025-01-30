package repository

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
	"github.com/stretchr/testify/assert"
)

// Тестирование слоя репозитория (RecomRepository) с моками вместо
// RelationalDataBase и KeyValueDatabse

// Моки:
type MockRelDB struct{}

func (m *MockRelDB) GetProductsByUserId(ctx context.Context, userId int) ([]int, error) {
	if userId == 4 {
		return nil, errors.New("some rel error")
	}
	return []int{1, 2, 2, 3}, nil
}
func (m *MockRelDB) AddProductUpdate(ctx context.Context, product *myproto.ProductAction) error {
	if product.Action == "ok" {
		return nil
	} else {
		return errors.New("some rel error")
	}

}
func (m *MockRelDB) AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) error {
	if user.UserInterests[0] == "ok" {
		return nil
	} else {
		return errors.New("some rel error")
	}
}

type MockKVDB struct{}

func (m *MockKVDB) GetRecom(ctx context.Context, userId int) ([]int, error) {
	if userId == 1 {
		return []int{1, 2, 3}, nil
	} else if userId == 2 {
		return nil, nil
	} else if userId == 3 {
		return nil, errors.New("some KV error")
	}
	return nil, nil
}

func (m *MockKVDB) SetRecom(ctx context.Context, userId int, productIds []int) error {
	if userId == 5 {
		return errors.New("some KV error")
	}
	return nil
}

func (m *MockKVDB) DelRecom(ctx context.Context, userId int) error {
	if userId == 6 {
		return errors.New("some KV error")
	}
	return nil
}

func (m *MockKVDB) DelAll(ctx context.Context) error {
	if ctx.Value([4]string{"code"}) == 1 {
		return errors.New("some KV error")
	}
	return nil
}

// Непосредственно тестирование
func TestRecomRepository_GetRecommendations_CorrectGetFromKV(t *testing.T) {

	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	products, err := r.GetRecommendations(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, products)
}

func TestRecomRepository_GetRecommendations_CorrectGetFromRelDB(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	products, err := r.GetRecommendations(context.Background(), 2)
	assert.NoError(t, err)
	assert.NotNil(t, products)
}

func TestRecomRepository_GetRecommendations_KVError(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	products, err := r.GetRecommendations(context.Background(), 3)
	assert.NoError(t, err)
	assert.NotNil(t, products)
}

func TestRecomRepository_GetRecommendations_RelError(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	products, err := r.GetRecommendations(context.Background(), 4)
	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestRecomRepository_GetRecomendations_KVErrorSet(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	products, err := r.GetRecommendations(context.Background(), 5)
	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestRecomRepository_AddProductUpdate_Correct(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddProductUpdate(context.Background(), &myproto.ProductAction{
		Action: "ok",
	})

	assert.NoError(t, err)
}

func TestRecomRepository_AddProductUpdate_Incorrect(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddProductUpdate(context.Background(), &myproto.ProductAction{
		Action: "error",
	})

	assert.Error(t, err)
}

func TestRecomRepository_AddProductUpdate_IncorrectErrorFromKvV(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	err := r.AddProductUpdate(context.WithValue(context.Background(), [4]string{"code"}, 1), &myproto.ProductAction{
		Action: "ok",
	})

	assert.Error(t, err)
}

func TestRecomRepository_AddUserUpdate_Correct(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddUserUpdate(context.Background(), &myproto.UserUpdate{
		UserInterests: []string{"ok"},
	})

	assert.NoError(t, err)
}

func TestRecomRepository_AddUserUpdate_Incorrect(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddUserUpdate(context.Background(), &myproto.UserUpdate{
		UserInterests: []string{"error"},
	})

	assert.Error(t, err)
}

func TestRecomRepository_AddUserUpdate_IncorrectKVError(t *testing.T) {
	r := NewRecomRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddUserUpdate(context.Background(), &myproto.UserUpdate{
		UserId:        6,
		UserInterests: []string{"ok"},
	})

	assert.Error(t, err)
}
