package repository

import (
	"fmt"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sqlx.DB
}

// установка соединения с базой, паника в случае ошибки
func NewPostgresDB(cfg config.DBConfig) *PostgresDB {

	db := sqlx.MustConnect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname, cfg.Sslmode))

	return &PostgresDB{
		db: db,
	}
}
