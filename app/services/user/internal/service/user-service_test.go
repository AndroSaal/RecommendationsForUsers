package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/stretchr/testify/assert"
)

//Тестирование сервисного слоя с моками на MailSender и repository.Repository

// Мок Сервсиа отправки писем
type MockMailSender struct{}

func (m *MockMailSender) SendMail(ctx context.Context, email string, code string) error {
	if email == "test@test.com" {
		return errors.New("Such email not exist")
	}
	return nil
}

// Мок Слоя репозиториев
type MockRepository struct{}

func (m *MockRepository) AddNewUser(ctx context.Context, user *entities.UserInfo, code string) (int, error) {
	if user.UsrDesc == "Incorrect User" {
		return 0, errors.New("IncorrectUser")
	}
	return 1, nil
}
func (m *MockRepository) GetUserById(ctx context.Context, id int) (*entities.UserInfo, error) {
	if id == 6 {
		return nil, errors.New("some repository level error")
	}
	return &entities.UserInfo{}, nil
}
func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error) {
	if email == "incorrect" {
		return nil, errors.New("Incorrect email")
	}
	return &entities.UserInfo{}, nil
}
func (m *MockRepository) VerifyCode(ctx context.Context, userId int, code string) (bool, error) {
	if code == "0" {
		return false, errors.New("Incorrect code")
	}
	return true, nil
}
func (m *MockRepository) UpdateUser(ctx context.Context, userId int, user *entities.UserInfo) error {
	if user.UsrDesc == "Incorrect User" {
		return errors.New("Incorrect User")
	}
	return nil
}

func TestUserService_CreateUser_Correct(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	userId, err := service.CreateUser(context.Background(), &entities.UserInfo{})

	assert.NoError(t, err)
	assert.Equal(t, 1, userId)
}

func TestUserService_CreateUser_Incorrect(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	userId, err := service.CreateUser(context.Background(), &entities.UserInfo{
		UsrDesc: "Incorrect User",
	})

	assert.Error(t, err)
	assert.Equal(t, 0, userId)
}

func TestUserService_CreateUser_NotExistingEmail(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	userId, err := service.CreateUser(context.Background(), &entities.UserInfo{
		Email: "test@test.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, userId)
}

func TestUserService_VerifyCode_Correct(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	isVerified, err := service.VerifyCode(context.Background(), 1, "21384")
	assert.NoError(t, err)
	assert.True(t, isVerified)
}

func TestUserService_VerifyCode_IncorrectCode(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	isVerified, err := service.VerifyCode(context.Background(), 1, "0")
	assert.Error(t, err)
	assert.False(t, isVerified)
}

func TestUserService_GetUserById_Correct(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	user, err := service.GetUserById(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserService_GetUserById_Incorrect(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	user, err := service.GetUserById(context.Background(), 6)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserService_GetUserByEmail_Correct(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	user, err := service.GetUserByEmail(context.Background(), "correct")

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserService_GetUserByEmail_Incorrect(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	user, err := service.GetUserByEmail(context.Background(), "incorrect")

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserService_UpdateUser_Correct(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := service.UpdateUser(context.Background(), 1, &entities.UserInfo{})

	assert.NoError(t, err)
}

func TestUserService_UpdateUser_Incorrect(t *testing.T) {
	//Создаем сервис
	service := NewUserService(&MockMailSender{}, &MockRepository{}, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	err := service.UpdateUser(context.Background(), 1, &entities.UserInfo{
		UsrDesc: "Incorrect User",
	})

	assert.Error(t, err)
}
