package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
)

// имплементация интерфейса Service
type UserService struct {
	repo repository.Repository // интерфейс для взаимодействия со слоем репозиториев
	log  *slog.Logger          // логгер для трейсов и логирования
	mail MailSender            // интерфейс для отправки писем
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
func (s *UserService) CreateUser(ctx context.Context, user *entities.UserInfo) (int, error) {
	fi := "internal.User.CreateUser" // используется для отслеживания ошибок

	// генерация кода
	code := generateCode()

	// добавление кода и пользователя в таблицы бд
	id, err := s.repo.AddNewUser(ctx, user, code)
	if err != nil {
		s.log.Debug("%s: Error adding new user: %v", fi, err)
		return 0, err
	}

	// отправка письма - опциональная функция,
	// если возникла ошибка - выполнение программы продолжится
	if err := s.mail.SendMail(ctx, user.Email, code); err != nil {
		s.log.Debug("%s: Error sending email: %v", fi, err)
	}

	return id, nil
}

// функция проверяет кода на валидность, если код совпадает с указанным в базе
// поле is_email_verified в таблице users меняется на true
func (s *UserService) VerifyCode(ctx context.Context, userId int, code string) (bool, error) {
	fi := "internal.User.VerifyCode"

	isVerified, err := s.repo.VerifyCode(ctx, userId, code)

	if err != nil {
		s.log.Debug(fmt.Sprintf("%s: %s", fi, err.Error()))
		return false, err
	}

	return isVerified, nil
}

// функция возвращает пользователя из базы по его id
func (s *UserService) GetUserById(ctx context.Context, id int) (*entities.UserInfo, error) {
	fi := "internal.User.GetUserById"

	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		s.log.Debug(fmt.Sprintf("%s: %s", fi, err.Error()))
		return nil, err

	}
	return user, nil
}

// функция возвращает пользователя из базы по его email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error) {
	fi := "internal.User.GetUserByEmail"

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Debug(fmt.Sprintf("%s: %s", fi, err.Error()))
		return nil, err
	}
	return user, nil
}

// функция заменяет информацию о пользователе в базе по его id
func (s *UserService) UpdateUser(ctx context.Context, userId int, user *entities.UserInfo) error {
	fi := "internal.User.UpdateUser"

	if err := s.repo.UpdateUser(ctx, userId, user); err != nil {
		s.log.Debug(fmt.Sprintf("%s: %s", fi, err.Error()))
		return err

	}
	return nil
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
