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
	UpdateUser(user *entities.UserInfo) error
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
	return r.relDB.AddNewUser(user, code)
}

func (r *UserRepository) GetUserById(id int) (*entities.UserInfo, error) {
	return r.relDB.GetUserById(id)
}

func (r *UserRepository) GetUserByEmail(email string) (*entities.UserInfo, error) {
	return r.relDB.GetUserByEmail(email)
}

func (r *UserRepository) VerifyCode(userId int, code string) (bool, error) {
	return r.relDB.VerifyCode(userId, code)
}

func (r *UserRepository) UpdateUser(user *entities.UserInfo) error {
	return r.relDB.UpdateUser(user)
}


