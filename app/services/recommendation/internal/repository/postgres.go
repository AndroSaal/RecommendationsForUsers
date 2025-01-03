package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RelationalDataBase interface {
	GetProductsByUserId(userId int) ([]int, error)
	AddProductUpdate(product *myproto.ProductAction) error
	AddUserUpdate(user *myproto.UserUpdate) error
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

// функция поиска id продуктов, в которых может быть заинтересован пользователь
func (p *PostgresDB) GetProductsByUserId(userId int) ([]int, error) {

	//начинаем транзакцию
	trx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	//проверка, существует ли пользователь с таким id
	rowCheck := trx.QueryRow(`SELECT id FROM users WHERE id = $1`, userId)
	if err := rowCheck.Scan(&userId); err != nil {
		trx.Rollback()
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return nil, err
	}

	//формируем запрос для получения всех kw_id пользователя (id его интересов)
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = $1`,
		kwIdField, userKwTable, userIdField,
	)

	//выполняем запрос
	rowsU, err := trx.Query(query, userId)
	if err != nil {
		trx.Rollback()
		return nil, err
	} else if err := rowsU.Err(); err != nil {
		trx.Rollback()
		return nil, err

	}
	defer rowsU.Close()

	//заполняем слайс интересами пользователя
	userInterests := make([]int, 0)
	for rowsU.Next() {
		var interestId int
		if err := rowsU.Scan(&interestId); err != nil {
			trx.Rollback()
			return nil, err
		}
		if interestId != 1 {
			userInterests = append(userInterests, interestId)
		} else {
			continue
		}
	}

	p.log.Info("User with id (%d) have interests with id %v", userId, userInterests)

	userRecommendations := make([]int, 0)
	for _, interestId := range userInterests {
		//формируем запрос для получения id продуктов, в которых есть ключевые слова пользователя
		query := fmt.Sprintf(
			`SELECT %s FROM %s WHERE %s = $1`,
			productIdField, productsKwTable, kwIdField,
		)
		//выполняем запрос
		rowsP, err := p.DB.Query(query, interestId)
		if err != nil {
			trx.Rollback()
			return nil, err
		}
		defer rowsP.Close()

		//заполняем слайс интересными пользователю продуктами
		for rowsP.Next() {
			var productId int
			if err := rowsP.Scan(&productId); err != nil {
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
	fi := "repository.AddUserUpdate"
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
		return fmt.Errorf("%s: %s %v", fi, query, err)
	}

	//Добавление ключевых слов (интересов) пользователя в таблицу keyWords и таблицу-связку
	if err := addKeyWords(int(user.UserId), tgx, user.UserInterests, userKwTable); err != nil {
		tgx.Rollback()
		p.log.Error("%s: Error adding User KeyWords (userId %d): %v", fi, user.UserId, err.Error(), err)
		return err
	} else {
		p.log.Info("%s: User KeyWords %v (userId %d) added", fi, user.UserInterests, user.UserId)
	}

	tgx.Commit()
	return nil
}

func (p *PostgresDB) AddProductUpdate(product *myproto.ProductAction) error {
	fi := "repository.AddProductUpdate"

	tgx, err := p.DB.Begin()
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

	//Добавление ключевых слов продукта в таблицу keyWords и таблицу-связку
	if err := addKeyWords(
		int(product.ProductId), tgx, product.ProductKeyWords, productsKwTable); err != nil {
		tgx.Rollback()
		p.log.Error("%s: Error adding Product KeyWords (productId %d): %s", fi, product.ProductId, err.Error(), err)
		return err
	} else {
		p.log.Info("%s: Product KeyWords %v (ProductId %d) added", fi, product.ProductKeyWords, product.ProductId)
	}

	tgx.Commit()
	return nil
}

// функция для добавления kw в таблицу keyWords и таблицу-связку userKw или productKw в зависимости от параметра table
func addKeyWords(id int, trx *sql.Tx, kw []string, table string) error {
	fi := "repository.addKeyWords"
	var idToinsert string

	if table == productsKwTable {
		idToinsert = productIdField
	} else if table == userKwTable {
		idToinsert = userIdField
	}
	//удаление страых связей
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		table, idToinsert,
	)
	if _, err := trx.Exec(query, id); err != nil {
		return fmt.Errorf("%s: %s %v", fi, query, err)
	}

	for _, keyWord := range kw {
		var keyWordId int
		//проверям есть ли такой keyword уже в таблице keyWords
		query := fmt.Sprintf(
			`SELECT %s FROM %s WHERE %s = $1`,
			idField, kwTable, kwNameField,
		)
		row := trx.QueryRow(query, keyWord)
		//получем id интереса
		if err := row.Scan(&keyWordId); err != nil {
			//	произошла ошибка - возвращаем ее
			if err != sql.ErrNoRows {
				return fmt.Errorf("%s: %s %v", fi, query, err)
			} else if keyWordId == 1 {
				continue
			}
			//	если нет в таблице такого keyWord - добавляем
			keyWordId, err = addKeyWord(trx, keyWord)
			if err != nil {
				return fmt.Errorf("%s: INSERT %v", fi, err)
			}
		}
		//добавление новых записей в таблицу связи
		querryKeyWordsProduct := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
			table,
			kwIdField, idToinsert,
		)
		if _, err := trx.Exec(querryKeyWordsProduct, keyWordId, id); err != nil {
			return fmt.Errorf("%s: %s %v", fi, querryKeyWordsProduct, err)
		}
	}
	return nil
}

func addKeyWord(trx *sql.Tx, keyWord string) (int, error) {
	var kwId int
	//формируем запрос для добавления новой записи в таблицу keyWords
	queryAddKeyWords := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
		kwTable,
		kwNameField, idField,
	)
	//выполняем запрос
	row := trx.QueryRow(queryAddKeyWords, keyWord)

	//получем id интереса
	if err := row.Scan(&kwId); err != nil {
		return 0, fmt.Errorf("can't get keyWord %s: %v", queryAddKeyWords, err)
	}
	return kwId, nil
}
