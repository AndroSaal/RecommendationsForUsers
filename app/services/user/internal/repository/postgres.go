package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/config"
	"github.com/jmoiron/sqlx"
	pq "github.com/lib/pq"
)

type RelationalDataBase interface {
	AddNewUser(user *entities.UserInfo, code string) (int, error)
	GetUserById(id int) (*entities.UserInfo, error)
	GetUserByEmail(email string) (*entities.UserInfo, error)
	VerifyCode(userId int, code string) (bool, error)
	UpdateUser(userId int, user *entities.UserInfo) error
}

// имплементация RelationalDataBase интерфейса
type PostgresDB struct {
	DB *sqlx.DB
}

// установка соединения с базой, паника в случае ошибки
func NewPostgresDB(cfg config.DBConfig) *PostgresDB {

	DB, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname, cfg.Sslmode))

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return &PostgresDB{
		DB: DB,
	}
}

func (p *PostgresDB) AddNewUser(user *entities.UserInfo, code string) (int, error) {

	var userId int
	//начинаем транзакцию
	trx, err := p.DB.Begin()
	if err != nil {
		return 0, err
	}

	//формируем запрос для добавления новой записи в таблицу users
	queryAddUser := fmt.Sprintf(
		`INSERT INTO %s 
		 (%s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING %s`,
		usersTable,
		emailPole, usernamePole, passwordPole, describtionPole, isEmailVerifiedPole, agePole,
		id,
	)
	//выполняем запрос по добавлению
	row := trx.QueryRow(queryAddUser,
		user.Email, user.Usrname, user.Password, user.UsrDesc, false, user.UsrAge)

	//вычитывает полученный id
	if err := row.Scan(&userId); err != nil {
		trx.Rollback()
		if err.(*pq.Error).Code == "23505" {
			return 0, ErrAlreadyExists
		} else {
			return 0, err
		}
	}
	//формируем запрос для добавления новой записи в таблицу codes
	queryAddCode := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
		codesTable,
		codePole, userIdPole,
	)
	//выполняем запрос
	if _, err = trx.Exec(queryAddCode, code, userId); err != nil {
		trx.Rollback()
		return 0, err
	}

	//добавление интересов пользователя
	if _, err = addUserInterests(user, trx, userId); err != nil {
		trx.Rollback()
		return 0, err
	}
	//ураа все получилось, коммит
	trx.Commit()

	return userId, nil
}

func (p *PostgresDB) GetUserById(userId int) (*entities.UserInfo, error) {
	var userDB UserInfoForDB

	tgx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	querySelectUser := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`,
		all, usersTable, id,
	)
	row := tgx.QueryRow(querySelectUser, userId)

	if err := row.Scan(
		&userDB.UsrId, &userDB.Email, &userDB.Usrname, &userDB.Password,
		&userDB.UsrDesc, &userDB.UsrAge, &userDB.IsEmailValid,
	); err != nil {
		tgx.Rollback()
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}

		return nil, err
	}

	queryOfAllInterests := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`,
		interestIdPole, userInterestsTable, userIdPole,
	)

	rows, err := tgx.Query(queryOfAllInterests, userId)
	if err != nil {
		tgx.Rollback()
		return nil, err
	}

	queryToSelectInterestName := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = $1`,
		intersestPole, interestsTable, id,
	)

	interests := make([]entities.UserInterest, 0)
	for rows.Next() {
		var (
			interest   entities.UserInterest
			interestId int
		)

		if err := rows.Scan(&interestId); err != nil {
			tgx.Rollback()
			return nil, err
		}

		interestRow := p.DB.QueryRow(queryToSelectInterestName, interestId)
		if err := interestRow.Scan(&interest); err != nil {
			tgx.Rollback()
			return nil, err
		}

		interests = append(interests, interest)
	}

	tgx.Commit()

	return &entities.UserInfo{
		UsrId:           userId,
		Usrname:         userDB.Usrname,
		Email:           userDB.Email,
		Password:        userDB.Password,
		UsrDesc:         entities.UserDiscription(userDB.UsrDesc),
		UserInterests:   interests,
		UsrAge:          entities.UserAge(userDB.UsrAge),
		IsEmailVerified: userDB.IsEmailValid,
	}, nil
}

func (p *PostgresDB) GetUserByEmail(email string) (*entities.UserInfo, error) {
	var userId int

	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`,
		id, usersTable, emailPole,
	)
	row := p.DB.QueryRow(query, email)

	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return nil, err
	}

	return p.GetUserById(userId)
}

func (p *PostgresDB) VerifyCode(userId int, code string) (bool, error) {
	var (
		codeFromDB string
	)

	tgx, err := p.DB.Begin()
	if err != nil {
		return false, err
	}

	//формирование запроса к базе
	querySelectCode := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`,
		codePole, codesTable, userIdPole,
	)

	//выполняем запрос,
	row := tgx.QueryRow(querySelectCode, userId)

	//получаем запись
	if err := row.Scan(&codeFromDB); err != nil {
		tgx.Rollback()
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return false, err
	}

	if codeFromDB == code {
		//формируем текст запроса
		queryToAddVerification := fmt.Sprintf(
			`UPDATE %s
			SET %s = true 
			WHERE %s = $1`,
			usersTable,
			isEmailVerifiedPole,
			id,
		)

		//выполняем запрос
		if _, err := tgx.Exec(queryToAddVerification, userId); err != nil {
			tgx.Rollback()
			return false, err
		}
	} else {
		tgx.Rollback()
		return false, nil
	}
	tgx.Commit()
	return true, nil
}

func (p *PostgresDB) UpdateUser(userId int, user *entities.UserInfo) error {

	tgx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	rowCheck := tgx.QueryRow(`SELECT id FROM users WHERE id = $1`, userId)

	//проверка что пользователь существует
	if err := rowCheck.Scan(&userId); err != nil {
		tgx.Rollback()
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return err

	}

	query := fmt.Sprintf(
		`UPDATE %s 
		 SET %s = $1, %s = $2, %s = $3, %s = $4, %s = $5 
		 WHERE %s = $6`,
		usersTable,
		usernamePole, emailPole, passwordPole, describtionPole, agePole,
		id,
	)

	if _, err := tgx.Exec(query,
		user.Usrname, user.Email, user.Password, user.UsrDesc, user.UsrAge, userId,
	); err != nil {
		tgx.Rollback()
		if err.(*pq.Error).Code == "23503" {
			err = ErrNotFound
		}
		return err
	}

	//удание старых интересов пользователя
	queryDeleteInterests := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		userInterestsTable, userIdPole,
	)
	if _, err := tgx.Exec(queryDeleteInterests, userId); err != nil {
		tgx.Rollback()
		return err

	}

	//добавление новых интересов пользователя
	if _, err := addUserInterests(user, tgx, userId); err != nil {
		tgx.Rollback()
		return err
	}

	tgx.Commit()
	return nil
}

func addUserInterests(user *entities.UserInfo, trx *sql.Tx, userId int) (int, error) {
	for _, interest := range user.UserInterests {

		var interestId int
		//формируем запрос для добавления новой записи в таблицу interests
		queryAddInterest := fmt.Sprintf(`INSERT INTO %s (%s) VALUES ($1) RETURNING %s`,
			interestsTable,
			intersestPole, id,
		)
		//выполняем запрос
		row := trx.QueryRow(queryAddInterest, interest)

		//получем id интереса
		if err := row.Scan(&interestId); err != nil {
			return 0, fmt.Errorf("can't get interest id: %v", err)
		}

		//формируем запрос для добавления новой записи в таблицу user_interests
		querryInterestAndUser := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
			userInterestsTable,
			userIdPole, interestIdPole,
		)
		//добавляем ид юзера и его интерес в таблицу user_interests
		if _, err := trx.Exec(querryInterestAndUser, userId, interestId); err != nil {
			return 0, err
		}
	}
	return userId, nil
}

// func (p *PostgresDB) Close() {
// 	p.DB.Close()
// }
