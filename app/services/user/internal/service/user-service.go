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
	repo repository.Repository
	log  *slog.Logger
	mail MailSender
}

func NewUserService(mail MailSender, repo repository.Repository, log *slog.Logger) *UserService {
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

	//генерация кода
	code := generateCode()

	//добавление кода и польщоватеоя в таблицы бд
	id, err := s.repo.AddNewUser(user, code)
	if err != nil {
		s.log.Debug("%s: Error adding new user: %v", fi, err)
		return 0, err
	}

	//отправка письма
	if err := s.mail.SendMail(user.Email, code); err != nil {
		s.log.Debug("%s: Error sending email: %v", fi, err)
		return 0, err
	}

	

	return id, nil
}

// функция проверяет кода на валидность, если код совпадает с указанным в базе
// поле is_email_verified в таблице users меняется на true
func (s *UserService) VerifyCode(userId int, code string) (bool, error) {
	fi := "internal.User.VerifyCode"

	isVerified, err := s.repo.VerifyCode(userId, code)

	if err != nil {
		s.log.Debug(fmt.Sprintf("%s: %s", fi, err.Error()))
		return false, err
	}

	return isVerified, nil
}

// функция возвращает пользователя из базы по его id
func (s *UserService) GetUserById(id int) (*entities.UserInfo, error) {
	return s.repo.GetUserById(id)
}

// функция возвращает пользователя из базы по его email
func (s *UserService) GetUserByEmail(email string) (*entities.UserInfo, error) {
	return s.repo.GetUserByEmail(email)
}

// функция заменяет информацию о пользователе в базе по его id
func (s *UserService) UpdateUser(userId int, user *entities.UserInfo) error {
	return s.repo.UpdateUser(userId, user)
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
