package repository

import (
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
)

type Repository interface {
	AddNewUser(user *entities.UserInfo, code string) (int, error)
	GetUserById(id int) (*entities.UserInfo, error)
	GetUserByEmail(email string) (*entities.UserInfo, error)
	VerifyCode(userId int, code string) (bool, error)
	UpdateUser(userId int, user *entities.UserInfo) error
}

// имплементация Repository интерфейса
type UserRepository struct {
	relDB RelationalDataBase
	log   *slog.Logger
}

// слой репощитория - взаимодействие с Базами данных
func NewUserRepository(db *PostgresDB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		relDB: db,
		log:   log,
	}
}

func (r *UserRepository) AddNewUser(user *entities.UserInfo, code string) (int, error) {
	fi := "repository.UserRepository.AddNewUser"

	userId, err := r.relDB.AddNewUser(user, code)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return 0, err

	}

	return userId, nil
}

func (r *UserRepository) GetUserById(userId int) (*entities.UserInfo, error) {
	fi := "repository.UserRepository.GetUserById"

	user, err := r.relDB.GetUserById(userId)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entities.UserInfo, error) {
	fi := "repository.UserRepository.GetUserByEmail"

	user, err := r.relDB.GetUserByEmail(email)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return nil, err

	}
	return user, nil
}

func (r *UserRepository) VerifyCode(userId int, code string) (bool, error) {
	fi := "repository.UserRepository.VerifyCode"

	isVerified, err := r.relDB.VerifyCode(userId, code)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return false, err
	}
	return isVerified, nil
}

func (r *UserRepository) UpdateUser(userId int, user *entities.UserInfo) error {
	fi := "repository.UserRepository.UpdateUser"

	err := r.relDB.UpdateUser(userId, user)
	if err != nil {
		r.log.Error(fi + ": " + err.Error())
		return err
	}

	return nil
}
