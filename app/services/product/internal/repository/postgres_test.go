package repository

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Тестирование методов работы с базой данных Postgre

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

func TestPostgreDB_AddNewProduct_Correct(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "кино",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{"режиссер", "актёр", "камера"},
	}

	// добавление тестового продукта в таблицу
	idSource, err1 := dbConn.AddNewProduct(context.Background(), &product)
	assert.NoError(t, err1)
	assert.NotEmpty(t, idSource)

	// проверка что продукт добавился в таблицу products и получил тот же id, который нам вернулся
	var idFromDb int
	err2 := dbConn.DB.QueryRow("SELECT id FROM products WHERE category = $1", product.Category).Scan(&idFromDb)
	assert.NoError(t, err2)
	assert.Equal(t, idSource, idFromDb)

	// проверка, что за тестовым продуктом закрпились указанные key-ворды
	keyWordsFromDB := make([]string, 0)
	rows, err3 := dbConn.DB.Query(
		"SELECT kw_name FROM keyWords a INNER JOIN (SELECT * FROM product_keyWord WHERE product_id = $1) b ON b.kw_id = a.id",
		idFromDb)
	assert.NoError(t, err3)
	for rows.Next() {
		var keyWord string
		err5 := rows.Scan(&keyWord)
		assert.NoError(t, err5)
		keyWordsFromDB = append(keyWordsFromDB, keyWord)
	}
	assert.Equal(t, product.ProductKeyWords, keyWordsFromDB)

}

func TestPostgreDB_AddNewProduct_Incorrect1(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{"режиссер", "актёр", "камера"},
	}

	// добавление тестового продукта в таблицу
	idSource, err1 := dbConn.AddNewProduct(context.Background(), &product)
	assert.Error(t, err1)
	assert.Empty(t, idSource)
}

func TestPostgreDB_AddNewProduct_Incorrect2(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "гонки",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{""},
	}

	// добавление тестового продукта в таблицу
	idSource, err1 := dbConn.AddNewProduct(context.Background(), &product)
	assert.Error(t, err1)
	assert.Empty(t, idSource)
}

func TestPostgreDB_UpdateProduct_Correct(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "кино",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{"режиссер", "актёр", "камера", "сценарий"},
	}

	// добавление тестового продукта в таблицу
	err1 := dbConn.UpdateProduct(context.Background(), 1, &product)
	assert.NoError(t, err1)

	// проверка что продукт добавился в таблицу products и получил тот же id, который нам вернулся
	var idFromDb, idSource int = 0, 1
	err2 := dbConn.DB.QueryRow("SELECT id FROM products WHERE category = $1", product.Category).Scan(&idFromDb)
	assert.NoError(t, err2)
	assert.Equal(t, idSource, idFromDb)

	// проверка, что за тестовым продуктом закрпились указанные key-ворды
	keyWordsFromDB := make([]string, 0)
	rows, err3 := dbConn.DB.Query(
		"SELECT kw_name FROM keyWords a INNER JOIN (SELECT * FROM product_keyWord WHERE product_id = $1) b ON b.kw_id = a.id ",
		idFromDb)
	assert.NoError(t, err3)
	for rows.Next() {
		var keyWord string
		err5 := rows.Scan(&keyWord)
		assert.NoError(t, err5)
		keyWordsFromDB = append(keyWordsFromDB, keyWord)
	}
	assert.Equal(t, product.ProductKeyWords, keyWordsFromDB)

}

func TestPostgreDB_UpdateProduct_Incorrect1(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "rbyj",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{"режиссер", "актёр", "камера"},
	}

	// добавление тестового продукта в таблицу
	err1 := dbConn.UpdateProduct(context.Background(), 100, &product)
	assert.Error(t, err1)
}

func TestPostgreDB_UpdateProduct_Incorrect2(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "гонки",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{""},
	}

	// добавление тестового продукта в таблицу
	err1 := dbConn.UpdateProduct(context.Background(), 1, &product)
	assert.Error(t, err1)
}

func TestPostgreDB_UpdateProduct_Incorrect3(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()
	// тестовый продукт
	product := entities.ProductInfo{
		Category:        "",
		Description:     "корректно",
		Status:          "availible",
		ProductKeyWords: []string{"режиссер", "актёр", "камера"},
	}

	// добавление тестового продукта в таблицу
	err1 := dbConn.UpdateProduct(context.Background(), 1, &product)
	assert.Error(t, err1)
}

func TestPostgre_DeleteProduct_Correct(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	err := dbConn.DeleteProduct(context.Background(), 1)
	assert.NoError(t, err)
}

func TestPostgre_DeleteProduct_Incorrect(t *testing.T) {

	cfg := loadConf()

	// коннект к бд (Маст)
	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))
	// закрываем коннект, выводим ошибку
	defer func() {
		if err := dbConn.DB.Close(); err != nil {
			t.Error(errors.New("Ошибка закрытия БД" + err.Error()))
		}
	}()

	err := dbConn.DeleteProduct(context.Background(), 100)
	assert.Error(t, err)
}
