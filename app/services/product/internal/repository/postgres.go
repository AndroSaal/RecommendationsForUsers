package repository

import (
	"database/sql"
	"fmt"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/pkg/config"
	"github.com/jmoiron/sqlx"
	pq "github.com/lib/pq"
)

type RelationalDataBase interface {
	AddNewProduct(product *entities.ProductInfo) (int, error)
	UpdateProduct(productId int, uproduct *entities.ProductInfo) error
	DeleteProduct(productId int) error
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

func (p *PostgresDB) AddNewProduct(product *entities.ProductInfo) (int, error) {

	var productId int

	//начинаем транзакцию
	trx, err := p.DB.Begin()
	if err != nil {
		return 0, err
	}

	//формируем запрос для добавления новой записи в таблицу users
	query := fmt.Sprintf(
		`INSERT INTO %s 
		 (%s, %s, %s) VALUES ($1, $2, $3) 
		 RETURNING %s`,
		productsTable,
		categoryField, describtionField, statusField,
		id,
	)

	//выполняем запрос по добавлению нового продукта
	row := p.DB.QueryRow(query,
		product.Category, product.Description, product.Status)

	//вычитываем полученный id
	if err := row.Scan(&productId); err != nil {
		trx.Rollback()
		return 0, err
	}

	//добавление ключевых слов продукта
	if err = addProductKeyWords(product, trx, productId); err != nil {
		trx.Rollback()
		return 0, err
	}
	//ураа все получилось, коммит
	trx.Commit()

	return productId, nil
}

func (p *PostgresDB) UpdateProduct(productId int, product *entities.ProductInfo) error {

	tgx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	//проверка что пользователь существует
	rowCheck := tgx.QueryRow(`SELECT id FROM products WHERE id = $1`, productId)
	if err := rowCheck.Scan(&productId); err != nil {
		tgx.Rollback()
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return err
	}

	query := fmt.Sprintf(
		`UPDATE %s 
		 SET %s = $1, %s = $2, %s = $3 
		 WHERE %s = $4`,
		productsTable,
		categoryField, describtionField, statusField,
		id,
	)

	if _, err := tgx.Exec(query,
		product.Category, product.Description, product.Status,
		productId,
	); err != nil {
		tgx.Rollback()
		if err.(*pq.Error).Code == "23503" {
			err = ErrNotFound
		}
		return err
	}

	//удание старых ключевых слов продукта
	queryDeleteInterests := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		productKwTable, productIdField,
	)
	if _, err := tgx.Exec(queryDeleteInterests, productId); err != nil {
		tgx.Rollback()
		return err

	}

	//добавление новых интересов пользователя
	if err := addProductKeyWords(product, tgx, productId); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

func (p *PostgresDB) DeleteProduct(productId int) error {

	tgx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	//проверка что пользователь существует
	rowCheck := p.DB.QueryRow(`SELECT id FROM products WHERE id = $1`, productId)
	if err := rowCheck.Scan(&productId); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		tgx.Rollback()
		return err
	}

	query := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		productsTable,
		id,
	)

	if _, err := p.DB.Exec(query, productId); err != nil {
		tgx.Rollback()
		return err
	}

	//удание старых ключевых слов продукта
	queryDeleteInterests := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		productKwTable, productIdField,
	)
	if _, err := tgx.Exec(queryDeleteInterests, productId); err != nil {
		tgx.Rollback()
		return err

	}

	return nil

}

func addProductKeyWords(product *entities.ProductInfo, trx *sql.Tx, productId int) error {
	for _, keyWord := range product.ProductKeyWords {

		var keyWordId int
		//формируем запрос для добавления новой записи в таблицу tags
		queryAddKeyWords := fmt.Sprintf(`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
			kwTable,
			kwNameField, id,
		)
		//выполняем запрос
		row := trx.QueryRow(queryAddKeyWords, keyWord)

		//получем id интереса
		if err := row.Scan(&keyWordId); err != nil {
			return fmt.Errorf("can't get keyWord id: %v", err)
		}

		//формируем запрос для добавления новой записи в таблицу products_tags
		querryKeyWordsProduct := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
			productKwTable,
			kwIdField, productIdField,
		)
		//добавляем ид юзера и его интерес в таблицу user_interests
		if _, err := trx.Exec(querryKeyWordsProduct, keyWordId, productId); err != nil {
			return err
		}
	}
	return nil
}
