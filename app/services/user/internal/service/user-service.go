package service

import "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"

// имплементация интерфейса Service
type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
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
