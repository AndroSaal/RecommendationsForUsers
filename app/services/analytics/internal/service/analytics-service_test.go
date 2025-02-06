package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"testing"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

type MockRepository struct{}

func (m *MockRepository) AddProductUpdate(ctx context.Context, product *myproto.ProductAction) error {
	if product.Action == "error" {
		return errors.New("some error")
	}
	return nil
}
func (m *MockRepository) AddUserUpdate(ctx context.Context, user *myproto.UserUpdate) error {
	if user.UserInterests[0] == "error" {
		return errors.New("some error")
	}
	return nil
}

func TestAnalyticsService_AddProductUpdate_Correct(t *testing.T) {
	service := NewAnalyticsService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	var message myproto.ProductAction = myproto.ProductAction{
		ProductId: 1,
		Action:    "add",
	}
	marshalledMessage, err := proto.Marshal(&message)
	if err != nil {
		t.Error(err)
	}

	err1 := service.AddProductData(context.Background(), &sarama.ConsumerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(marshalledMessage),
	})

	assert.NoError(t, err1)
}

func TestAnalyticsService_AddProductUpdate_Incorrect(t *testing.T) {
	service := NewAnalyticsService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	var message myproto.ProductAction = myproto.ProductAction{
		ProductId: 1,
		Action:    "add",
	}
	marshalledMessage, err := json.Marshal(&message)
	if err != nil {
		t.Error(err)
	}

	err1 := service.AddProductData(context.Background(), &sarama.ConsumerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(marshalledMessage),
	})

	assert.Error(t, err1)
}

func TestAnalyticsService_AddProductUpdate_IncorrectErrorONRepoLevel(t *testing.T) {
	service := NewAnalyticsService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	var message myproto.ProductAction = myproto.ProductAction{
		ProductId: 1,
		Action:    "error",
	}
	marshalledMessage, err := proto.Marshal(&message)
	if err != nil {
		t.Error(err)
	}

	err1 := service.AddProductData(context.Background(), &sarama.ConsumerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(marshalledMessage),
	})

	assert.Error(t, err1)
}

func TestAnalyticsService_AddUserUpdate_Correct(t *testing.T) {
	service := NewAnalyticsService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	var message myproto.UserUpdate = myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"add"},
	}

	marshalledMessage, err := proto.Marshal(&message)
	if err != nil {
		t.Error(err)
	}

	err1 := service.AddUserData(context.Background(), &sarama.ConsumerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(marshalledMessage),
	})

	assert.NoError(t, err1)
}

func TestAnalyticsService_AddUserUpdate_Incorrect(t *testing.T) {
	service := NewAnalyticsService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	var message myproto.UserUpdate = myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"add"},
	}

	marshalledMessage, err := json.Marshal(&message)
	if err != nil {
		t.Error(err)
	}

	err1 := service.AddUserData(context.Background(), &sarama.ConsumerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(marshalledMessage),
	})

	assert.Error(t, err1)
}

func TestAnalyticsService_AddUserUpdate_IncorrectErrorONRepoLevel(t *testing.T) {
	service := NewAnalyticsService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	var message myproto.UserUpdate = myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"error"},
	}
	marshalledMessage, err := proto.Marshal(&message)
	if err != nil {
		t.Error(err)
	}

	err1 := service.AddUserData(context.Background(), &sarama.ConsumerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(marshalledMessage),
	})

	assert.Error(t, err1)
}
