package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/config"
	"github.com/jmoiron/sqlx"
)

type RelationalDataBase interface {
	AddProductUpdate(product *myproto.ProductAction) error
	AddUserUpdate(user *myproto.UserUpdate) error
}

// имплементация RelationalDataBase интерфейса
type PostgresDB struct {
	DB *sqlx.DB
}

// установка соединения с базой, паника в случае ошиби
func NewPostgresDB(cfg config.DBConfig) *PostgresDB {

	db := sqlx.MustConnect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname, cfg.Sslmode))

	return &PostgresDB{
		DB: db,
	}
}

func (p *PostgresDB) AddUserUpdate(user *myproto.UserUpdate) error {

	tgx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	//добавление id пользователя в табблицу users
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) ON CONFLICT DO NOTHING`,
		usersTable, idField,
	)
	if _, err := tgx.Exec(query, user.UserId); err != nil {
		tgx.Rollback()
		return err
	}
	//Добавление информации о обновлении
	query = fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
		userTimestampsTable, userIdField, timestampField,
	)
	row := tgx.QueryRow(query, user.UserId)

	var timestamp time.Time
	if err := row.Scan(&timestamp); err != nil {
		tgx.Rollback()
		return err
	}

	//Добавление ключевых слов (интересов) пользователя в таблицу user_upadates
	if err := addKeyWords(
		timestamp, tgx, user.UserInterests, userUpdatesTable); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

func (p *PostgresDB) AddProductUpdate(product *myproto.ProductAction) error {

	tgx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	//добавление id пользователя в табблицу users
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) ON CONFLICT DO NOTHING`,
		productsTable, idField,
	)
	if _, err := tgx.Exec(query, product.ProductId); err != nil {
		tgx.Rollback()
		return err
	}

	//Добавление информации о обновлении
	query = fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
		productsTimestampsTable, productIdField, timestampField,
	)
	row := tgx.QueryRow(query, product.ProductId)

	var timestamp time.Time
	if err := row.Scan(&timestamp); err != nil {
		tgx.Rollback()
		return err
	}

	//Добавление ключевых слов (интересов) пользователя в таблицу user_upadates
	if err := addKeyWords(
		timestamp, tgx, product.ProductKeyWords, productUpsetesTable); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

func addKeyWords(timestamp time.Time, trx *sql.Tx, kw []string, table string) error {
	var idToinsert string

	if table == userUpdatesTable {
		idToinsert = kwField
	} else if table == productUpsetesTable {
		idToinsert = userInterestsField
	}

	str := strings.Join(kw, ",")
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		table, timestampField, idToinsert,
	)
	if _, err := trx.Exec(query, timestamp, str); err != nil {
		return err
	}
	return nil
}
