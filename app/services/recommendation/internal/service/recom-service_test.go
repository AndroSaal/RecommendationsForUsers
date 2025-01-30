package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"testing"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

//тесты для слоя сервиса
// type RecommendationService struct {
// 	repo repository.Repository
// 	log  *slog.Logger
// }

type MockRepository struct{}

func (m *MockRepository) GetRecommendations(ctx context.Context, userId int) ([]int, error) {
	if userId == 1 {
		return nil, errors.New("some error")
	} else if userId == 2 {
		return []int{1, 2, 3}, nil
	} else {
		return []int{}, nil
	}
}

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

func TestRecommendationService_GetRecommendations_Correct(t *testing.T) {
	service := NewRecommendationService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	res, err := service.GetRecommendations(context.Background(), 2)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestRecommendationService_GetRecommendations_Incorrect(t *testing.T) {
	service := NewRecommendationService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	res, err := service.GetRecommendations(context.Background(), 1)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestRecommendationService_GetRecommendations_CorrectButNoRecommendation(t *testing.T) {
	service := NewRecommendationService(
		&MockRepository{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	res, err := service.GetRecommendations(context.Background(), 3)

	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestRecommendationService_AddProductUpdate_Correct(t *testing.T) {
	service := NewRecommendationService(
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

func TestRecommendationService_AddProductUpdate_Incorrect(t *testing.T) {
	service := NewRecommendationService(
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

func TestRecommendationService_AddProductUpdate_IncorrectErrorONRepoLevel(t *testing.T) {
	service := NewRecommendationService(
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

func TestRecommendationService_AddUserUpdate_Correct(t *testing.T) {
	service := NewRecommendationService(
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

func TestRecommendationService_AddUserUpdate_Incorrect(t *testing.T) {
	service := NewRecommendationService(
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

func TestRecommendationService_AddUserUpdate_IncorrectErrorONRepoLevel(t *testing.T) {
	service := NewRecommendationService(
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
