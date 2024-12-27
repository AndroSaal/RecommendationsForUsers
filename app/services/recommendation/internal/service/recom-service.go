package service

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/repository"
	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

// имплементация интерфейса Service
type RecommendationService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewRecommendationService(repo repository.Repository, log *slog.Logger) *RecommendationService {
	return &RecommendationService{
		repo: repo,
		log:  log,
	}
}

// функция вызывает метод репозитория по добавлению нового продукта
func (s *RecommendationService) GetRecommendations(userId int) ([]int, error) {
	return s.repo.GetRecommendations(userId)
}

func (s *RecommendationService) AddProductData(msg *sarama.ConsumerMessage) error {
	fi := "service.RecommendationService.AddProductData"
	var (
		product *myproto.ProductAction
	)

	//из слайса байт в структуру
	if err := proto.Unmarshal(msg.Value, product); err != nil {
		s.log.Error(fi, ": ", "Error unmarshaling product entity: ", err)
		return err
	}

	//отправляем структуру в бд
	if err := s.repo.AddProductUpdate(product); err != nil {
		s.log.Error(fi, ": ", "Error adding product entity: ", err)
		return err
	}

	return nil
}

func (s *RecommendationService) AddUserData(msg *sarama.ConsumerMessage) error {
	fi := "service.RecommendationService.AddUserData"
	var (
		user *myproto.UserUpdate
	)

	//из слайса байт в структуру
	if err := proto.Unmarshal(msg.Value, user); err != nil {
		s.log.Error(fi, ": ", "Error unmarshaling user entity: ", err)
		return err
	}

	//отправляем структуру в бд
	if err := s.repo.AddUserUpdate(user); err != nil {
		s.log.Error(fi, ": ", "Error adding user entity: ", err)
		return err
	}

	return nil
}
