package repository

import "log/slog"

type Repository interface{}

type UserRepository struct {
	relDB *PostgresDB
	log   *slog.Logger
}

// слой репощитория - взаимодействие с Базами данных
func NewUserRepository(db *PostgresDB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		relDB: db,
		log:   log,
	}
}
