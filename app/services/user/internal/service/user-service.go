package service

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
)

// имплементация интерфейса Service
type UserService struct {
	repo *repository.Repository
	log  *slog.Logger
	mail *Mail
}

func NewUserService(mail *Mail, repo *repository.Repository, log *slog.Logger) *UserService {
	return &UserService{
		mail: mail,
		repo: repo,
		log:  log,
	}
}

// функция вызывает метод репозитория по добавлению нового пользователя
// если не возвращается ошибка -> отправляется письмо на указанную почту
func (s *UserService) CreateUser(user *entities.UserInfo) (int, error) {
	fi := "internal.User.CreateUser"

	id, err := s.repo.AddNewUser(user, code)
	if err != nil {
		s.log.Debug("%s: Error adding new user: %v", fi, err)
		return 0, err
	}

	code := generateCode()

	if err := s.mail.SendMail(user.Email, code); err != nil {
		s.log.Debug("%s: Error sending email: %v", fi, err)
		return 0, err
	}

	return id, nil
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
	fi := "internal.User.VerifyCode"

	isVerified, err := s.repo.GetCodeFromEmail(userId, code)

	if err != nil {
		s.log.Debug(fmt.Sprintf("%s: %s", fi, err.Error()))
		return false, err
	}

	return isVerified, nil
}

func generateCode() string {

	rand.Seed(time.Now().UnixNano())

	// Генерация четырех случайных чисел от 0 до 9
	code := []int{
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
	}

	// Соединение чисел в одну строку
	finalcode := fmt.Sprintf("%d%d%d%d%d", code[0], code[1], code[2], code[3], code[4])

	return finalcode

}
