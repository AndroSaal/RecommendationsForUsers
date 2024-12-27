package repository

import (
	"database/sql"
	"fmt"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/jmoiron/sqlx"
	pq "github.com/lib/pq"
)

type RelationalDataBase interface {
	GetProductsByUserId(userId int) ([]int, error)
	AddProductUpdate(product *myproto.ProductAction) error
	AddUserUpdate(user *myproto.UserUpdate) error
}

// имплементация RelationalDataBase интерфейса
type PostgresDB struct {
	db *sqlx.DB
}

// установка соединения с базой, паника в случае ошиби
func NewPostgresDB(cfg config.DBConfig) *PostgresDB {

	db := sqlx.MustConnect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname, cfg.Sslmode))

	return &PostgresDB{
		db: db,
	}
}

// функция поиска id продуктов, в которых может быть заинтересован пользователь
func (p *PostgresDB) GetProductsByUserId(userId int) ([]int, error) {

	//начинаем транзакцию
	trx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	//формируем запрос для получения всех kw_id пользователя (id его интересов)
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = $1`,
		kwIdField, userKwTable, userIdField,
	)

	//выполняем запрос
	rows, err := p.db.Query(query, userId)
	if err != nil {
		trx.Rollback()
		return nil, err
	}
	defer rows.Close()

	//заполняем слайс интересами пользователя
	userInterests := make([]int, 3)
	for rows.Next() {
		var interestId int
		if err := rows.Scan(&interestId); err != nil {
			trx.Rollback()
			return nil, err
		}
		userInterests = append(userInterests, interestId)
	}

	userRecommendations := make([]int, 3)
	for _, interestId := range userInterests {
		//формируем запрос для получения id продуктов, в которых есть ключевые слова пользователя
		query := fmt.Sprintf(
			`SELECT %s FROM %s WHERE %s = $1`,
			productIdField, productsKwTable, kwIdField,
		)
		//выполняем запрос
		rows, err := p.db.Query(query, interestId)
		if err != nil {
			trx.Rollback()
			return nil, err
		}
		defer rows.Close()

		//заполняем слайс интересными пользователю продуктами
		for rows.Next() {
			var productId int
			if err := rows.Scan(&productId); err != nil {
				trx.Rollback()
				return nil, err
			}
			userRecommendations = append(userRecommendations, productId)
		}
	}

	trx.Commit()

	return userRecommendations, nil
}

func (p *PostgresDB) AddUserUpdate(user *myproto.UserUpdate) error {

	tgx, err := p.db.Begin()
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

	//Добавление ключевых слов (интересов) пользователя в таблицу keyWords и таблицу-связку
	if err := addKeyWords(int(user.UserId), tgx, user.UserInterests, userKwTable); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

func (p *PostgresDB) AddProductUpdate(product *myproto.ProductAction) error {

	tgx, err := p.db.Begin()
	if err != nil {
		return err
	}

	//добавление id продукта в таблицу products
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) ON CONFLICT DO NOTHING`,
		productsTable, idField,
	)
	if _, err := tgx.Exec(query, product.ProductId); err != nil {
		tgx.Rollback()
		return err
	}

	//Добавление ключевых слов (интересов) пользователя в таблицу keyWords и таблицу-связку
	if err := addKeyWords(
		int(product.ProductId), tgx, product.ProductKeyWords, productsKwTable); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

// функция для добавления kw в таблицу keyWords и таблицу-связку userKw или productKw в зависимости от параметра table
func addKeyWords(id int, trx *sql.Tx, kw []string, table string) error {
	var idToinsert string

	if table == productsKwTable {
		idToinsert = productIdField
	} else if table == usersTable {
		idToinsert = userIdField
	}
	//удаление страых связей
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1 IF EXIST`,
		table, idToinsert,
	)
	if _, err := trx.Exec(query, id); err != nil {
		return err
	}

	for _, keyWord := range kw {
		var keyWordId int
		//формируем запрос для добавления новой записи в таблицу keyWords
		queryAddKeyWords := fmt.Sprintf(
			`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
			kwTable,
			kwNameField, idField,
		)
		//выполняем запрос
		row := trx.QueryRow(queryAddKeyWords, keyWord)

		//получем id интереса
		if err := row.Scan(&keyWordId); err != nil {
			//ошибка нарушения уникальности - интерес с таким именем уже есть
			if err.(*pq.Error).Code == "23503" {
				//ищем id существующего интереса
				query := fmt.Sprintf(
					`SELECT %s FROM %s WHERE %s = $1`,
					idField, kwTable, kwNameField,
				)
				row := trx.QueryRow(query, keyWord)
				if err := row.Scan(&keyWordId); err != nil {
					return fmt.Errorf("can't get keyWord id even after select :-( : %v", err)
				}
			} else {
				return fmt.Errorf("can't get keyWord id: %v", err)
			}
		}

		//добавление новых записей в таблицу связи
		querryKeyWordsProduct := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
			table,
			kwIdField, idToinsert,
		)
		if _, err := trx.Exec(querryKeyWordsProduct, keyWordId, id); err != nil {
			return err
		}
	}
	return nil
}
