package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Тесты для RelationalDatabase, проверяющие как отрабатывают запросы к базе
func loadConf() config.DBConfig {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	opt := os.Getenv("DB_OPT_DC")
	dbname := os.Getenv("DB_NAME")

	if user == "" {
		log.Fatal("loadConfig - Can't get user env")
	}

	if password == "" {
		log.Fatal("loadConfig - Can't get password env")
	}

	if host == "" {
		log.Fatal("loadConfig - Can't get host env")
	}

	if port == "" {
		log.Fatal("loadConfig - Can't get port env")
	}

	if dbname == "" {
		log.Fatal("loadConfig - Can't get dbname env")
	}

	if opt == "" {
		log.Fatal("loadConfig - Can't get opt env")
	}

	return config.DBConfig{
		Username: user,
		Password: password,
		Host:     host,
		Port:     port,
		Sslmode:  opt,
		Dbname:   dbname,
	}
}
func TestPostgresDB_AddNewUser_Correct(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	// тестовый юзер для вставки в таблицу
	var user *entities.UserInfo = &entities.UserInfo{
		Usrname:       "test",
		Email:         "test_test@test.com",
		Password:      "test",
		UsrDesc:       "more test words test test",
		UserInterests: []entities.UserInterest{"test1", "test2", "test3"},
		UsrAge:        15,
	}

	// добавление тестового юзера в таблицу
	idSource, err1 := dbConn.AddNewUser(context.Background(), user, "5436")
	assert.NoError(t, err1)
	assert.NotEmpty(t, idSource)

	// проверка что юзер добавился в таблицу users и получил тот же id, который нам вернулся
	var idFromDb int
	err2 := dbConn.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&idFromDb)
	assert.NoError(t, err2)
	assert.Equal(t, idSource, idFromDb)

	// проверка, что за тестовым юзером закрпились указанные интересы
	interestsFromDB := make(entities.UserInterests, 0)
	rows, err3 := dbConn.DB.Query(
		"SELECT interest FROM interests a INNER JOIN (SELECT * FROM user_interests WHERE user_id = $1) b ON b.interest_id = a.id ",
		idFromDb)
	assert.NoError(t, err3)
	for rows.Next() {
		var interest entities.UserInterest
		err5 := rows.Scan(&interest)
		assert.NoError(t, err5)
		interestsFromDB = append(interestsFromDB, interest)
	}
	assert.Equal(t, user.UserInterests, interestsFromDB)

	// проверка что в таблицу codes добавился нужный код
	var codeFromDB string
	err4 := dbConn.DB.QueryRow("SELECT email_code FROM codes WHERE user_id = $1", idSource).Scan(&codeFromDB)
	assert.NoError(t, err4)
	assert.Equal(t, "5436", codeFromDB)
}

func TestPostgresDB_AddNewUser_CorrectAlreadyExist(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	// тестовый юзер для вставки в таблицу
	var user *entities.UserInfo = &entities.UserInfo{
		Usrname:       "test",
		Email:         "test_test@test.com",
		Password:      "test",
		UsrDesc:       "more test words test test",
		UserInterests: []entities.UserInterest{"test1", "test2", "test3"},
		UsrAge:        15,
	}

	// добавление тестового юзера в таблицу
	idSource, err1 := dbConn.AddNewUser(context.Background(), user, "5436")

	// проверка что возникает ошибка, т.к. пользователь с таким email уже существует
	assert.Error(t, err1)
	assert.Empty(t, idSource)
}

func TestPostgresDB_GetUserById_CorrectId(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	// вставленный в таблицы ранее юзер
	var userSource *entities.UserInfo = &entities.UserInfo{
		UsrId:         1,
		Usrname:       "test",
		Email:         "test_test@test.com",
		Password:      "test",
		UsrDesc:       "more test words test test",
		UserInterests: []entities.UserInterest{"test1", "test2", "test3"},
		UsrAge:        15,
	}

	user, err := dbConn.GetUserById(context.Background(), 1)

	// проверка что получаем юзера, добавленного в предыдущем тесте
	assert.NoError(t, err)
	assert.Equal(t, userSource, user)
}

func TestPostgresDB_GetUserById_InсorrectId(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	user, err := dbConn.GetUserById(context.Background(), 100)

	//проверка, что пользователя с таким id нет
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestPostgresDB_GetUserByEmail_CorrectEmail(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	var userSource *entities.UserInfo = &entities.UserInfo{
		UsrId:         1,
		Usrname:       "test",
		Email:         "test_test@test.com",
		Password:      "test",
		UsrDesc:       "more test words test test",
		UserInterests: []entities.UserInterest{"test1", "test2", "test3"},
		UsrAge:        15,
	}

	user, err := dbConn.GetUserByEmail(context.Background(), "test_test@test.com")

	//проверка, что пользователь с таким email есть
	assert.NoError(t, err)
	assert.Equal(t, userSource, user)
}

func TestPostgresDB_GetUserByEmail_IncorrectEmail(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	user, err := dbConn.GetUserByEmail(context.Background(), "test@test.com")

	//проверка, что пользователь с таким email есть
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestPostgresDB_VerifyCode_CorrectCode(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	isVerified, err := dbConn.VerifyCode(context.Background(), 1, "5436")
	assert.NoError(t, err)
	assert.True(t, isVerified)
}

func TestPostgresDB_VerifyCode_IncorrectCode(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	isVerified, err := dbConn.VerifyCode(context.Background(), 1, "5488")
	assert.NoError(t, err)
	assert.False(t, isVerified)
}

func TestPostgresDB_UpdateUser_Correct(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	var user *entities.UserInfo = &entities.UserInfo{
		Usrname:         "test",
		Email:           "test_test@test.com",
		Password:        "test",
		UsrDesc:         "more test words test test",
		UserInterests:   []entities.UserInterest{"test1", "test2"},
		UsrAge:          16,
		IsEmailVerified: true,
	}

	// Представим, что пользователь повзрослел на год и из интересов пропал test3
	err1 := dbConn.UpdateUser(context.Background(), 1, user)
	assert.NoError(t, err1)

	//Берём из таблицы этого юзера и проверяем действительно ли он обновился
	userFromDB, err2 := dbConn.GetUserByEmail(context.Background(), "test_test@test.com")
	assert.NoError(t, err2)
	user.UsrId = 1
	assert.Equal(t, user, userFromDB)
}

func TestPostgresDB_UpdateUser_IncorrectId(t *testing.T) {
	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg)
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	var user *entities.UserInfo = &entities.UserInfo{
		Usrname:       "test",
		Email:         "test_test@test.com",
		Password:      "test",
		UsrDesc:       "more test words test test",
		UserInterests: []entities.UserInterest{"test1"},
		UsrAge:        17,
	}

	//пробуем обновить юзера с несуществующим id
	err1 := dbConn.UpdateUser(context.Background(), 100, user)

	//проверем, что возвращается ошибка
	assert.Error(t, err1)

}
