package repository

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Тестирование методов работы реляционной БД

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

func TestPostgreDB_AddUserUpdate_Correct(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"машины", "рыба", "анчоусы"},
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен без ошибок
	assert.NoError(t, err)

	// проверка, что  функция добавила всё как задумывалось
	rows, err := dbConn.DB.Query(fmt.Sprintf(
		"SELECT kw_id FROM user_kw WHERE user_id = %d", user.UserId,
	))

	assert.NoError(t, err)
	// Инициализируем срез с нулевой длиной
	interestsId := make([]int, 0)

	for rows.Next() {
		var interest int
		// Сканируем значение в переменную interest
		err = rows.Scan(&interest)
		assert.NoError(t, err)

		// Добавляем значение в срез
		interestsId = append(interestsId, interest)
	}

	// Не забудьте закрыть rows после использования
	rows.Close()
	assert.Equal(t, 3, len(interestsId))

	for i, elem := range user.UserInterests {
		row := dbConn.DB.QueryRow(fmt.Sprintf(
			"SELECT kw_name FROM keyWords WHERE id = %d", interestsId[i],
		))
		kwName := ""
		err := row.Scan(&kwName)
		assert.NoError(t, err)
		assert.Equal(t, elem, kwName)
	}
}

func TestPostgreDB_AddUserUpdate_CorrectAlreadyExist(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"машины", "рыба"},
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен c ошибкой
	assert.NoError(t, err)

	// проверка, что  функция добавила всё как задумывалось
	rows, err := dbConn.DB.Query(fmt.Sprintf(
		"SELECT kw_id FROM user_kw WHERE user_id = %d", user.UserId,
	))

	assert.NoError(t, err)
	// Инициализируем срез с нулевой длиной
	interestsId := make([]int, 0)

	for rows.Next() {
		var interest int
		// Сканируем значение в переменную interest
		err = rows.Scan(&interest)
		assert.NoError(t, err)

		// Добавляем значение в срез
		interestsId = append(interestsId, interest)
	}

	// Не забудьте закрыть rows после использования
	rows.Close()
	assert.Equal(t, 2, len(interestsId))

	for i, elem := range user.UserInterests {
		row := dbConn.DB.QueryRow(fmt.Sprintf(
			"SELECT kw_name FROM keyWords WHERE id = %d", interestsId[i],
		))
		kwName := ""
		err := row.Scan(&kwName)
		assert.NoError(t, err)
		assert.Equal(t, elem, kwName)
	}
}

func TestPostgreDB_AddUserUpdate_Incorrect1(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        3,
		UserInterests: []string{},
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_AddUserUpdate_Incorrect2(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        -3,
		UserInterests: []string{},
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_AddUserUpdate_Incorrect3(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        4,
		UserInterests: nil,
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_AddUserUpdate_Incorrect4(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        4,
		UserInterests: []string{""},
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_AddProductUpdate_Correct(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var product myproto.ProductAction = myproto.ProductAction{
		ProductId:       1,
		ProductKeyWords: []string{"машины", "рыба", "анчоусы"},
	}

	err := dbConn.AddProductUpdate(context.Background(), &product)

	// проверка, что запрос выполнен без ошибок
	assert.NoError(t, err)

	// проверка, что  функция добавила всё как задумывалось
	rows, err := dbConn.DB.Query(fmt.Sprintf(
		"SELECT kw_id FROM product_kw WHERE product_id = %d", product.ProductId,
	))

	assert.NoError(t, err)
	// Инициализируем срез с нулевой длиной
	interestsId := make([]int, 0)

	for rows.Next() {
		var interest int
		// Сканируем значение в переменную interest
		err = rows.Scan(&interest)
		assert.NoError(t, err)

		// Добавляем значение в срез
		interestsId = append(interestsId, interest)
	}

	// Не забудьте закрыть rows после использования
	rows.Close()
	assert.Equal(t, 3, len(interestsId))

	for i, elem := range product.ProductKeyWords {
		row := dbConn.DB.QueryRow(fmt.Sprintf(
			"SELECT kw_name FROM keyWords WHERE id = %d", interestsId[i],
		))
		kwName := ""
		err := row.Scan(&kwName)
		assert.NoError(t, err)
		assert.Equal(t, elem, kwName)
	}
}

func TestPostgreDB_AddProductUpdate_CorrectAlreadyExist(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого продукта
	var product myproto.ProductAction = myproto.ProductAction{
		ProductId:       1,
		ProductKeyWords: []string{"машины", "рыба"},
	}

	err := dbConn.AddProductUpdate(context.Background(), &product)

	// проверка, что запрос выполнен без ошибок
	assert.NoError(t, err)

	// проверка, что  функция добавила всё как задумывалось
	rows, err := dbConn.DB.Query(fmt.Sprintf(
		"SELECT kw_id FROM product_kw WHERE product_id = %d", product.ProductId,
	))

	assert.NoError(t, err)
	// Инициализируем срез с нулевой длиной
	interestsId := make([]int, 0)

	for rows.Next() {
		var interest int
		// Сканируем значение в переменную interest
		err = rows.Scan(&interest)
		assert.NoError(t, err)

		// Добавляем значение в срез
		interestsId = append(interestsId, interest)
	}

	// Не забудьте закрыть rows после использования
	rows.Close()
	assert.Equal(t, 2, len(interestsId))

	for i, elem := range product.ProductKeyWords {
		row := dbConn.DB.QueryRow(fmt.Sprintf(
			"SELECT kw_name FROM keyWords WHERE id = %d", interestsId[i],
		))
		kwName := ""
		err := row.Scan(&kwName)
		assert.NoError(t, err)
		assert.Equal(t, elem, kwName)
	}
}

func TestPostgreDB_AddProductUpdate_Incorrect2(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var product myproto.ProductAction = myproto.ProductAction{
		ProductId:       -3,
		ProductKeyWords: []string{},
	}

	err := dbConn.AddProductUpdate(context.Background(), &product)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_AddProductUpdate_Incorrect3(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var product myproto.ProductAction = myproto.ProductAction{
		ProductId:       3,
		ProductKeyWords: []string{},
	}

	err := dbConn.AddProductUpdate(context.Background(), &product)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_AddProductUpdate_Incorrect4(t *testing.T) {

	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)

	}()

	// сущность добавляемого пользователя
	var product myproto.ProductAction = myproto.ProductAction{
		ProductId:       3,
		ProductKeyWords: []string{""},
	}

	err := dbConn.AddProductUpdate(context.Background(), &product)

	// проверка, что запрос выполнен без ошибок
	assert.Error(t, err)
}

func TestPostgreDB_GetProductsByUserId_Correct(t *testing.T) {
	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)
	}()

	recom, err := dbConn.GetProductsByUserId(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recom))

}

func TestPostgreDB_GetProductsByUserId_Incorrect(t *testing.T) {
	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)
	}()

	recom, err := dbConn.GetProductsByUserId(context.Background(), 100)
	assert.Error(t, err)
	assert.Nil(t, recom)

}

func TestPostgreDB_GetProductsByUserId_CorrectButNoRecom(t *testing.T) {
	// Подключение к Базе данных
	dbConn := NewPostgresDB(
		loadConf(), slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	// отключение
	defer func() {
		err := dbConn.DB.Close()
		assert.NoError(t, err)
	}()

	// сущность добавляемого пользователя
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        2,
		UserInterests: []string{"mega"},
	}

	err := dbConn.AddUserUpdate(context.Background(), &user)

	// проверка, что запрос выполнен без ошибок
	assert.NoError(t, err)

	recom, err := dbConn.GetProductsByUserId(context.Background(), 2)
	assert.NoError(t, err)
	assert.Equal(t, recom, []int{})

}
