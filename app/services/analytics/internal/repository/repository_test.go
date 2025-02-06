package repository

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/stretchr/testify/assert"
)

type MockRelDB struct{}

func (m *MockRelDB) AddProductUpdate(ctx context.Context, product *myproto.ProductAction) (time.Time, error) {
	if product.Action == "ok" || product.Action == "kv error" {
		return time.Now(), nil
	} else {
		return time.Now(), errors.New("some rel error")
	}
}

func (m *MockRelDB) AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) (time.Time, error) {
	if user.UserId == 0 {
		return time.Now(), nil
	} else {
		return time.Now(), errors.New("some rel error")
	}

}

type MockKVDB struct{}

func (m *MockKVDB) SetUserUpdate(ctx context.Context, user *myproto.UserUpdate, time time.Time) error {
	if user.UserInterests[0] == "ok" {
		return nil
	}
	return errors.New("some KV error")

}

func (m *MockKVDB) SetProductUpdate(ctx context.Context, product *myproto.ProductAction, time time.Time) error {
	if product.Action == "ok" {
		return nil
	}
	return errors.New("some KV error")
}

func (m *MockKVDB) Del(ctx context.Context, userId int) error {
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

func TestRecomRepository_AddProductUpdate_Correct(t *testing.T) {
	r := NewAnalyticsRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddProductUpdate(context.Background(), &myproto.ProductAction{
		Action: "ok",
	})

	assert.NoError(t, err)
}

func TestRecomRepository_AddProductUpdate_Incorrect(t *testing.T) {
	r := NewAnalyticsRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddProductUpdate(context.Background(), &myproto.ProductAction{
		Action: "error",
	})

	assert.Error(t, err)
}

func TestRecomRepository_AddProductUpdate_IncorrectErrorFromKvV(t *testing.T) {
	r := NewAnalyticsRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	err := r.AddProductUpdate(context.WithValue(context.Background(), [4]string{"code"}, 1), &myproto.ProductAction{
		Action: "kv error",
	})

	assert.Error(t, err)
}

func TestRecomRepository_AddUserUpdate_Correct(t *testing.T) {
	r := NewAnalyticsRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddUserUpdate(context.Background(), &myproto.UserUpdate{
		UserInterests: []string{"ok"},
	})

	assert.NoError(t, err)
}

func TestRecomRepository_AddUserUpdate_Incorrect(t *testing.T) {
	r := NewAnalyticsRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddUserUpdate(context.Background(), &myproto.UserUpdate{
		UserInterests: []string{"error"},
		UserId:        2,
	})

	assert.Error(t, err)
}

func TestRecomRepository_AddUserUpdate_IncorrectKVError(t *testing.T) {
	r := NewAnalyticsRepository(&MockRelDB{}, &MockKVDB{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := r.AddUserUpdate(context.Background(), &myproto.UserUpdate{
		UserId:        0,
		UserInterests: []string{"kv"},
	})

	assert.Error(t, err)
}
