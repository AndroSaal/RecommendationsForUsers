package repository

import (
	"context"
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
)

type Repository interface {
	AddNewUser(ctx context.Context, user *entities.UserInfo, code string) (int, error)
	GetUserById(ctx context.Context, id int) (*entities.UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error)
	VerifyCode(ctx context.Context, userId int, code string) (bool, error)
	UpdateUser(ctx context.Context, userId int, user *entities.UserInfo) error
}

// имплементация Repository интерфейса
type UserRepository struct {
	relDB RelationalDataBase
	log   *slog.Logger
}

// слой репозитория - взаимодействие с Базами данных
func NewUserRepository(db RelationalDataBase, log *slog.Logger) *UserRepository {
	return &UserRepository{
		relDB: db,
		log:   log,
	}
}

func (r *UserRepository) AddNewUser(ctx context.Context, user *entities.UserInfo, code string) (int, error) {
	fi := "repository.UserRepository.AddNewUser"

	userId, err := r.relDB.AddNewUser(ctx, user, code)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return 0, err

	}

	return userId, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId int) (*entities.UserInfo, error) {
	fi := "repository.UserRepository.GetUserById"

	user, err := r.relDB.GetUserById(ctx, userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error) {
	fi := "repository.UserRepository.GetUserByEmail"

	user, err := r.relDB.GetUserByEmail(ctx, email)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err

	}
	return user, nil
}

func (r *UserRepository) VerifyCode(ctx context.Context, userId int, code string) (bool, error) {
	fi := "repository.UserRepository.VerifyCode"

	isVerified, err := r.relDB.VerifyCode(ctx, userId, code)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return false, err
	}
	return isVerified, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userId int, user *entities.UserInfo) error {
	fi := "repository.UserRepository.UpdateUser"

	err := r.relDB.UpdateUser(ctx, userId, user)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}
