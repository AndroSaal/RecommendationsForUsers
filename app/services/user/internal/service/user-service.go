package service

import (
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
)

// имплементация интерфейса Service
type UserService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user entities.UserInfo) (int, error) {
	return 0, nil
}

func (s *UserService) GetUserById(id int) (entities.UserInfo, error) {
	return entities.UserInfo{}, nil
}

func (s *UserService) GetUserByEmail(email string) (entities.UserInfo, error) {
	return entities.UserInfo{}, nil
}

func (s *UserService) UpdateUser(user entities.UserInfo) error {
	return nil
}

func (s *UserService) SendMail(email *string) error {
	return nil
}

func (s *UserService) VerifyCode(userId int, code string) (bool, error) {
	return false, nil
}
