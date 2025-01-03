package repository

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RelationalDataBase interface {
	AddProductUpdate(product *myproto.ProductAction) (time.Time, error)
	AddUserUpdate(user *myproto.UserUpdate) (time.Time, error)
}

// имплементация RelationalDataBase интерфейса
type PostgresDB struct {
	DB  *sqlx.DB
	log *slog.Logger
}

// установка соединения с базой, паника в случае ошиби
func NewPostgresDB(cfg config.DBConfig, log *slog.Logger) *PostgresDB {

	db := sqlx.MustConnect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname, cfg.Sslmode))

	return &PostgresDB{
		DB:  db,
		log: log,
	}
}

func (p *PostgresDB) AddUserUpdate(user *myproto.UserUpdate) (time.Time, error) {
	fi := "repository.postgresDB.AddUserUpdate"

	tgx, err := p.DB.Begin()
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", fi, err)
	}

	//добавление id пользователя в табблицу users
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) ON CONFLICT DO NOTHING`,
		usersTable, idField,
	)
	if _, err := tgx.Exec(query, user.UserId); err != nil {
		tgx.Rollback()
		return time.Time{}, fmt.Errorf("%s: %w", fi, err)
	}
	//Добавление информации о обновлении
	query = fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES ($1, $2) RETURNING %s`,
		userUpdatesTable, userIdField, userInterestsField, timestampField,
	)
	str := strings.Join(user.UserInterests, ",")
	row := tgx.QueryRow(query, user.UserId, str)

	var timestamp time.Time
	if err := row.Scan(&timestamp); err != nil {
		tgx.Rollback()
		return time.Time{}, fmt.Errorf("%s: %w", fi, err)
	}

	p.log.Info(fmt.Sprintf("%s: SUCCESS added user update at %s", fi, timestamp))
	tgx.Commit()
	return timestamp, nil
}

func (p *PostgresDB) AddProductUpdate(product *myproto.ProductAction) (time.Time, error) {
	fi := "repository.postgresDB.AddProductUpdate"

	tgx, err := p.DB.Begin()
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", fi, err)
	}

	//добавление id пользователя в табблицу users
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) ON CONFLICT DO NOTHING`,
		productsTable, idField,
	)
	if _, err := tgx.Exec(query, product.ProductId); err != nil {
		tgx.Rollback()
		return time.Time{}, fmt.Errorf("%s: %w", fi, err)
	}

	//Добавление информации о обновлении
	query = fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES ($1, $2) RETURNING %s`,
		productUpdatesTable, productIdField, kwField, timestampField,
	)
	var str string
	if product.Action == "delete" {
		str = "DELETED"
	} else {
		str = strings.Join(product.ProductKeyWords, ",")
	}
	row := tgx.QueryRow(query, product.ProductId, str)

	var timestamp time.Time
	if err := row.Scan(&timestamp); err != nil {
		tgx.Rollback()
		return time.Time{}, fmt.Errorf("%s: %w", fi, err)
	}
	p.log.Info(fmt.Sprintf("%s: SUCCESS added product update at %s", fi, timestamp))

	tgx.Commit()
	return timestamp, nil
}
