package repository

import (
	"fmt"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RelationalDataBase interface {
	AddNewUser(user *entities.UserInfo, code string) (int, error)
	GetUserById(id int) (*entities.UserInfo, error)
	GetUserByEmail(email string) (*entities.UserInfo, error)
	VerifyCode(userId int, code string) (bool, error)
	UpdateUser(user *entities.UserInfo) error
}

// имплементация RelationalDataBase интерфейса
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

func (p *PostgresDB) AddNewUser(user *entities.UserInfo, code string) (int, error) {
	query := `INSERT INTO users (email, password, code) VALUES ($1, $2, $3)`
	return 0, nil
}

func (p *PostgresDB) GetUserById(id int) (*entities.UserInfo, error) {
	return nil, nil
}

func (p *PostgresDB) GetUserByEmail(email string) (*entities.UserInfo, error) {
	return nil, nil
}

func (p *PostgresDB) VerifyCode(userId int, code string) (bool, error) {
	return false, nil
}

func (p *PostgresDB) UpdateUser(user *entities.UserInfo) error {
	return nil
}
