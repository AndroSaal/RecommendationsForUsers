package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// тестирование транспортного слоя с моками на Sevice и Kafka

// мок сервиса
type MockService struct{}

func (m *MockService) CreateUser(ctx context.Context, usrInfo *entities.UserInfo) (int, error) {
	if usrInfo.UsrDesc == "Такой пользователь уже существует" { //симуляция юзер уже существует
		return 0, repository.ErrAlreadyExists
	} else if usrInfo.UsrDesc == "Такой пользователь вызовет проблемы на сервере" { //ошибка сервера
		return 0, errors.New("внутренняя ошибка сервера")
	}
	return 1, nil
}

func (m *MockService) GetUserById(ctx context.Context, id int) (*entities.UserInfo, error) {
	if id == 404 { //симуляция юзер не найден
		return nil, repository.ErrNotFound
	} else if id == 500 { //симуляция ошибка сервера
		return nil, errors.New("внутренняя ошибка сервера")
	}
	return nil, nil
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error) {
	if email == "user404@test.com" { //симуляция юзер не найден
		return nil, repository.ErrNotFound
	} else if email == "user500@test.com" { //симуляция ошибка сервера
		return nil, errors.New("внутренняя ошибка сервера")
	}
	return nil, nil
}

func (m *MockService) UpdateUser(ctx context.Context, id int, usrInfo *entities.UserInfo) error {
	if usrInfo.UsrDesc == "User Not Found" { //симуляция юзер уже существует
		return repository.ErrNotFound
	} else if usrInfo.UsrDesc == "Internal error" { //ошибка сервера
		return errors.New("внутренняя ошибка сервера")
	}
	return nil
}

func (m *MockService) VerifyCode(ctx context.Context, userId int, code string) (bool, error) {
	if userId == 1 {
		return true, nil
	} else if userId == 2 {
		return false, errors.New("внутренняя ошибка сервера")
	} else if userId == 3 {
		return false, repository.ErrNotFound
	}
	return false, nil
}

func NewMockService() *MockService {
	return &MockService{}
}

// мок кафки
type MockKafka struct{}

func (m *MockKafka) SendMessage(usrInfo entities.UserInfo) error {
	time.Sleep(5 * time.Microsecond)
	return nil
}

func (m *MockKafka) Close() error {
	return nil
}

func NewMockKafka() *MockKafka {
	return &MockKafka{}
}

// непосредственно тестирование

func TestUserHandler_InitRoutes_Correct(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)

	handler := NewHandler(
		NewMockService(),
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		NewMockKafka(),
	)

	// инициализируем маршруты
	router := handler.InitRoutes()

	// проверяем что маршруты инициализировались
	assert.NotNil(t, router)
}

// Функция регистрауции нового пользователя
// Корректный запрос - никаких ошибок
func TestUserHandler_SingUpUser_CorrectRequest(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTest",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	//проверяем ответ
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	var response map[string]int
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1, response["userId"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

// Некорректный запрос - отсутствие имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_MissingUsername(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	userInfo := UserInfo{
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'UserInfo.Usrname' Error:Field validation for 'Usrname' failed on the 'required' tag",
		response["reason"])
}

// Некорректный запрос - отсутствие почты пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_MissingEmail(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	userInfo := UserInfo{
		Username:    "testTest",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'UserInfo.Email' Error:Field validation for 'Email' failed on the 'required' tag",
		response["reason"])
}

// Некорректный запрос - отсутствие описания пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_MissingDescription(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username  string   `json:"username"`
		Email     string   `json:"email"`
		Password  string   `json:"password"`
		Interests []string `json:"interests"`
		Age       int      `json:"age"`
	}

	userInfo := UserInfo{
		Username:  "testTest",
		Password:  "test0071",
		Email:     "testTest@test.com",
		Interests: []string{"FirstTest", "Second Test"},
		Age:       15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'UserInfo.UsrDesc' Error:Field validation for 'UsrDesc' failed on the 'required' tag",
		response["reason"])
}

// Некорректный запрос - отсутствие интересов пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_MissingInterests(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Description string `json:"description"`
		Age         int    `json:"age"`
	}

	userInfo := UserInfo{
		Username:    "testTest",
		Password:    "test0071",
		Email:       "testTest@test.com",
		Description: "test test and test",
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'UserInfo.UserInterests' Error:Field validation for 'UserInterests' failed on the 'required' tag",
		response["reason"])
}

// Некорректный запрос - отсутсвие возраста пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_MissingAge(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
	}

	userInfo := UserInfo{
		Username:    "testTest",
		Password:    "test0071",
		Email:       "testTest@test.com",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'UserInfo.UsrAge' Error:Field validation for 'UsrAge' failed on the 'required' tag",
		response["reason"])
}

// Некорректный запрос - ошибка валидации нарушение верхней границы количества символов
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_UsernameTooLong(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTestTTTTTTTTTTTEEEEEEEEESSSSSSSSTTTTTT",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid username: too long, max length is 20", response["reason"])
}

// Некорректный запрос - ошибка валидации нарушение нижней границы количества символов
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_UsernameTooShort(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "te",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid username: too short, min length is 3", response["reason"])
}

// Некорректный запрос - ошибка валидации регулярного выражения
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_UsernameRegexpErr(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "teJUH(фффф)",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, fmt.Sprintf(
		"invalid username: %s does not match regexp", userInfo.Username,
	), response["reason"])
}

// Некорректный запрос - ошибка валидации нарушение верхней границы количества символов
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_EmailTooLong(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "testTestTTTTTTTTTTTEEEEEEEEESSSSSSSSTTTTTTtesttesttesttetstetstetstes@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid email: too long, max length is 64", response["reason"])
}

// Некорректный запрос - ошибка валидации нарушение нижней границы количества символов
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_emailTooShort(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid email: too short, min length is 10", response["reason"])
}

// Некорректный запрос - ошибка валидации регулярного выражения
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_EmailRegexpErr(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "testtest.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid email: does not math regexp", response["reason"])
}

func TestUserHandler_SingUpUser_IncorrectRequest_PasswordTooLong(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "test@test.com",
		Password:    "testejciuwheicuhweicjosaknscdijweiojcoejndcowien0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid password: too long, max length is 32", response["reason"])
}

// Некорректный запрос - ошибка валидации нарушение нижней границы количества символов
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_PasswordTooShort(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "testtest@test.com",
		Password:    "test",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid password: too short, min length is 8", response["reason"])
}

// Некорректный запрос - ошибка валидации регулярного выражения
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_PasswordRegexpErr(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "test@test.com",
		Password:    "@testtest",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid password: does not math regexp", response["reason"])
}

func TestUserHandler_SingUpUser_IncorrectRequest_DescriptionRegexpErr(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username: "test",
		Email:    "test@test.com",
		Password: "testtest",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
			"Sed do eiusmod tempor incididunt ut labore et dolore magna " +
			"aliqua. Ut enim ad minim veniam, quis nostrud exercitation " +
			"ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis" +
			"aute irure dolor in reprehenderit in voluptate velit esse cillum " +
			"dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat " +
			"non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." +
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
			"Sed do eiusmod tempor incididunt ut labore et dolore magna " +
			"aliqua. Ut enim ad minim veniam, quis nostrud exercitation " +
			"ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis" +
			"aute irure dolor in reprehenderit in voluptate velit esse cillum " +
			"dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat " +
			"non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." +
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
			"Sed do eiusmod tempor incididunt ut labore et dolore magna " +
			"aliqua. Ut enim ad minim veniam, quis nostrud exercitation " +
			"ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis" +
			"aute irure dolor in reprehenderit in voluptate velit esse cillum ",
		Interests: []string{"FirstTest", "Second Test"},
		Age:       15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user discription: too long, max length is 1024", response["reason"])
}

// Некорректный запрос - ошибка валидации нарушение нижней границы количества символов
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_InterestsTooShort(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "testtest@test.com",
		Password:    "testaksn",
		Description: "test test and test",
		Interests:   []string{"F", "Second Test"},
		Age:         15,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "user interests[0]: invalid user intersest: too short, min length is 3", response["reason"])
}

// Некорректный запрос - ошибка валидации регулярного выражения
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_AgeTooOld(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "test@test.com",
		Password:    "testet@testtest",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         151,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user age: must be between 5 and 150", response["reason"])
}

// Некорректный запрос - ошибка валидации регулярного выражения
// для имени пользователя
func TestUserHandler_SingUpUser_IncorrectRequest_AgeTooYoung(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "test@test.com",
		Password:    "testet@testtest",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         4,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user age: must be between 5 and 150", response["reason"])
}

// Код оошибки 409
func TestUserHandler_SingUpUser_CorrectRequestButAlreadyExist(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "test@test.com",
		Password:    "testet@testtest",
		Description: "Такой пользователь уже существует",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         18,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "user with such email already exists", response["reason"])
}

func TestUserHandler_SingUpUser_CorrectRequestButInternalServerError(t *testing.T) {

	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запрос
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "test",
		Email:       "test@test.com",
		Password:    "testet@testtest",
		Description: "Такой пользователь вызовет проблемы на сервере",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         18,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/sign-up", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// отправляем запрос
	handler.signUpUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "внутренняя ошибка сервера", response["reason"])
}

// Тестирование получения пользователя по его ID
func TestUserHandler_GetUserById_Correct(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/userId?userId=1", nil)

	handler.getUserById(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestUserHandler_GetUserById_Incorrect_WrongId(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/userId?userId=тесто", nil)

	handler.getUserById(c)

	//проверяем
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "strconv.Atoi: parsing \"тесто\": invalid syntax", response["reason"])
}

func TestUserHandler_GetUserById_Incorrect_No(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.getUserById(c)

	//проверяем
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "userId parametr does not exist in path", response["reason"])
}

func TestUserHandler_GetUserById_Incorrect_ValidateErrorLess0(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/userId?userId=-1", nil)

	handler.getUserById(c)

	//проверяем
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user id: can`t be less or equel 0", response["reason"])
}

func TestUserHandler_GetUserById_Incorrect_ValidateErrorEq0(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/userId?userId=0", nil)

	handler.getUserById(c)

	//проверяем
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user id: can`t be less or equel 0", response["reason"])
}

func TestUserHandler_GetUserById_CorrectButNotFound(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/userId?userId=404", nil)

	handler.getUserById(c)

	//проверяем
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "user not found", response["reason"])
}

func TestUserHandler_GetUserById_CorrectButInternalError(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/userId?userId=500", nil)

	handler.getUserById(c)

	//проверяем
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "внутренняя ошибка сервера", response["reason"])
}

// Тестирование получения пользователя по его email
func TestUserHandler_GetUserByEmail_Correct(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/email?email=test@test.com", nil)

	handler.getUserByEmail(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestUserHandler_GetUserByEmail_IncorrectValidationEmail(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/email?email=testtest.com", nil)

	handler.getUserByEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid email: does not math regexp", response["reason"])
}

func TestUserHandler_GetUserByEmail_CorrectButNotFound(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/email?email=user404@test.com", nil)

	handler.getUserByEmail(c)

	//проверяем
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "user not found", response["reason"])
}

func TestUserHandler_GetUserByEmail_CorrectButInternalError(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/user/sign-up/email?email=user500@test.com", nil)

	handler.getUserByEmail(c)

	//проверяем
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "внутренняя ошибка сервера", response["reason"])
}

func TestUserHandler_GetUserByEmail_Incorrect_No(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.getUserByEmail(c)

	//проверяем
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "email parametr does not exist in path", response["reason"])
}

func TestUserHandler_EditUser_Correct(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTest",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestUserHandler_EditUser_IncorrectNoParametr(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTest",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)

	handler.editUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "userId parametr does not exist in path", response["reason"])
}

func TestUserHandler_EditUser_IncorrectNoBody(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid request", response["reason"])
}

func TestUserHandler_EditUser_IncorrectNoOneOfField(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTest",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'UserInfo.Email' Error:Field validation for 'Email' failed on the 'required' tag", response["reason"])
}

func TestUserHandler_EditUser_IncorrectValidationBody(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "t",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid username: too short, min length is 3", response["reason"])
}

func TestUserHandler_EditUser_IncorrectValidationParam(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "t",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "0"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user id: can`t be less or equel 0", response["reason"])
}

func TestUserHandler_EditUser_IncorrectParam(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "t",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "test test and test",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "Kot"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "strconv.Atoi: parsing \"Kot\": invalid syntax", response["reason"])
}

func TestUserHandler_EditUser_CorrectButNotFound(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTest",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "User Not Found",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, repository.ErrNotFound.Error(), reponse["reason"])
}

func TestUserHandler_EditUser_CorrectButInternalErorr(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// структура для тела запроса
	type UserInfo struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
		Age         int      `json:"age"`
	}

	// тело тестового запроса
	userInfo := UserInfo{
		Username:    "testTest",
		Email:       "test@test.com",
		Password:    "test0071",
		Description: "Internal error",
		Interests:   []string{"FirstTest", "Second Test"},
		Age:         15,
	}

	json_data, err := json.Marshal(userInfo)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(json_data)

	// Исправляем путь и добавляем параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/edit", reader)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.editUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "внутренняя ошибка сервера", reponse["reason"])
}

func TestUserHandler_VerifyEmail_Correct(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=80744", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	var reponse map[string]bool
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, true, reponse["verified"])
}

func TestUserHandler_VerifyEmail_IncorrectMissingUserId(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=80744", nil)

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "userId parametr does not exist in path", reponse["reason"])
}

func TestUserHandler_VerifyEmail_IncorrectMissingEmailCode(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "code parametr does not exist in query", reponse["reason"])
}

func TestUserHandler_VerifyEmail_IncorrectValidationCode(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=8074", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "invalid code: does not match regexp", reponse["reason"])
}

func TestUserHandler_VerifyEmail_IncorrectUserId(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=80744", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "Maxon"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "strconv.Atoi: parsing \"Maxon\": invalid syntax", reponse["reason"])
}

func TestUserHandler_VerifyEmail_IncorrectValidationUserId(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=80744", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "0"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "invalid user id: can`t be less or equel 0", reponse["reason"])
}

func TestUserHandler_VerifyEmail_CorrectButInternalServerError(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=80744", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "2"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, "внутренняя ошибка сервера", reponse["reason"])
}

func TestUserHandler_VerifyEmail_CorrectButNotFound(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)
	handler := &UserHandler{
		service: NewMockService(),
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		kafka: NewMockKafka(),
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// путь и параметр userId
	c.Request, _ = http.NewRequest("PATCH", "/user/sign-up/1/verify-email?code=80744", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "3"},
	}

	handler.verifyEmail(c)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var reponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &reponse)
	assert.Equal(t, repository.ErrNotFound.Error(), reponse["reason"])
}
