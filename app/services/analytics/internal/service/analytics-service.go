package service

import (
	"context"
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/repository"
	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

// имплементация интерфейса Service
type AnalyticsService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewAnalyticsService(repo repository.Repository, log *slog.Logger) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
		log:  log,
	}
}

func (s *AnalyticsService) AddProductData(ctx context.Context, msg *sarama.ConsumerMessage) error {
	fi := "service.AnalyticsServiceAnalyticsService.AddProductData"
	var product myproto.ProductAction

	//из слайса байт в структуру
	if err := proto.Unmarshal(msg.Value, &product); err != nil {
		s.log.Error(fi, ": ", "Error unmarshaling product entity: ", err.Error(), err)
		return err
	}

	//отправляем структуру в бд
	if err := s.repo.AddProductUpdate(ctx, &product); err != nil {
		s.log.Error(fi, ": ", "Error adding product entity: ", err.Error(), err)
		return err
	}

	return nil
}

func (s *AnalyticsService) AddUserData(ctx context.Context, msg *sarama.ConsumerMessage) error {
	fi := "service.AnalyticsService.AddUserData"
	var user myproto.UserUpdate

	//из слайса байт в структуру
	if err := proto.Unmarshal(msg.Value, &user); err != nil {
		s.log.Error(fi, ": ", "Error unmarshaling user entity: ", err.Error(), err)
		return err
	}

	//отправляем структуру в бд
	if err := s.repo.AddUserUpdate(ctx, &user); err != nil {
		s.log.Error(fi, ": ", "Error adding user entity: ", err.Error(), err)
		return err
	}

	return nil
}
