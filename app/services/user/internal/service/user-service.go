package service

import "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"

// имплементация интерфейса Service
type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user entities.UserInfo) (entities.UserId, error) {
	return 0, nil
}

func (s *UserService) GetUserById(id entities.UserId) (entities.UserInfo, error) {
	return entities.UserInfo{}, nil
}

func (s *UserService) GetUserByEmail(email entities.Email) (entities.UserInfo, error) {
	return entities.UserInfo{}, nil
}

func (s *UserService) UpdateUser(user entities.UserInfo) error {
	return nil
}

func (s *UserService) SendMail(email *entities.Email) error {
	return nil
}
