package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RelationalDataBase interface {
	AddNewProduct(ctx context.Context, product *entities.ProductInfo) (int, error)
	UpdateProduct(ctx context.Context, productId int, uproduct *entities.ProductInfo) error
	DeleteProduct(ctx context.Context, productId int) error
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

func (p *PostgresDB) AddNewProduct(ctx context.Context, product *entities.ProductInfo) (int, error) {

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
		p.log.Info(fmt.Sprintf("error scanning productID: %s", err.Error()))
		return 0, err
	}

	//добавление ключевых слов продукта
	if err = addProductKeyWords(product, trx, productId, p.log); err != nil {
		trx.Rollback()
		p.log.Info(fmt.Sprintf("error while addProductKeyWords: %s", err.Error()))
		return 0, err
	}
	//ураа все получилось, коммит
	trx.Commit()

	return productId, nil
}

func (p *PostgresDB) UpdateProduct(ctx context.Context, productId int, product *entities.ProductInfo) error {

	tgx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	//проверка что продукт существует
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
	if err := addProductKeyWords(product, tgx, productId, p.log); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

func (p *PostgresDB) DeleteProduct(ctx context.Context, productId int) error {

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

func addProductKeyWords(product *entities.ProductInfo, trx *sql.Tx, productId int, log *slog.Logger) error {
	for _, keyWord := range product.ProductKeyWords {

		var keyWordId int
		// проверяем, есть ли уже такой интерес в таблице keyWords
		queryGetKeyWord := fmt.Sprintf(`SELECT id FROM %s WHERE %s = $1`, kwTable, kwNameField)

		rowSelect := trx.QueryRow(queryGetKeyWord, keyWord)

		if err := rowSelect.Scan(&keyWordId); errors.Is(err, sql.ErrNoRows) { //если такого нет keyWord
			log.Info(fmt.Sprintf("postgres: keyWord %s not found, will be add", keyWord))
			//формируем запрос для добавления новой записи в таблицу keyWords
			queryAddKeyWord := fmt.Sprintf(`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
				kwTable,
				kwNameField, id,
			)
			//выполняем запрос
			rowInsert := trx.QueryRow(queryAddKeyWord, keyWord)

			//получем id интереса
			if err := rowInsert.Scan(&keyWordId); err != nil {
				return fmt.Errorf("can't get keyWord id: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("can't get keyWord id: %v", err)
		}

		//формируем запрос для добавления новой записи в таблицу product_keyWord
		querryKeyWordsProduct := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
			productKwTable,
			kwIdField, productIdField,
		)
		//добавляем ид продукта и его key-words в таблицу связку
		if _, err := trx.Exec(querryKeyWordsProduct, keyWordId, productId); err != nil {
			return err
		}
	}
	return nil
}
