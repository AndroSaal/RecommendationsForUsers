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
	fi := "service.RecommendationService.GetRecommendations"
	productIds, err := s.repo.GetRecommendations(userId)
	if err != nil {
		s.log.Error("%s: Error trying get product ids: %v", fi, err)
		return nil, err
	}
	s.log.Info("%s: Got request about user with id %d", fi, userId)
	if len(productIds) == 0 {
		productIds = nil
	}

	return productIds, nil
}

func (s *RecommendationService) AddProductData(msg *sarama.ConsumerMessage) error {
	fi := "service.RecommendationService.AddProductData"
	var product myproto.ProductAction

	//из слайса байт в структуру
	if err := proto.Unmarshal(msg.Value, &product); err != nil {
		s.log.Error("%s: Error trying Unmarshal product data: %v", fi, err)
		return err
	}

	s.log.Info("%s: Got user id %d, user interests %v", fi, product.ProductId, product.ProductKeyWords)

	//отправляем структуру в бд
	if err := s.repo.AddProductUpdate(&product); err != nil {
		s.log.Error("%s: Error trying add product data: %v", fi, err)
		return err
	}

	return nil
}

func (s *RecommendationService) AddUserData(msg *sarama.ConsumerMessage) error {
	fi := "service.RecommendationService.AddUserData"
	var (
		user myproto.UserUpdate
	)

	//из слайса байт в структуру
	if err := proto.Unmarshal(msg.Value, &user); err != nil {
		s.log.Error("%s: Error trying Unmarshal user data: %v", fi, err)
		return err
	} else {
		s.log.Info("%s: Unmarshal user data: id: %d, Interests %v", fi, user.UserId, user.UserInterests)
	}

	s.log.Info("%s: Got user id %d, user interests %v", fi, user.UserId, user.UserInterests)

	//отправляем структуру в бд
	if err := s.repo.AddUserUpdate(&user); err != nil {
		s.log.Error("%s: Error trying add user data: %v", fi, err)
		return err
	}

	return nil
}
