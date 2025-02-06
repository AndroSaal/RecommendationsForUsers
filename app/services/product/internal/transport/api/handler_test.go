package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

// Мок для стурктуры сервиса
type MockService struct{}

func (m *MockService) CreateProduct(ctx context.Context, user *entities.ProductInfo) (int, error) {
	if user.Description == "интернал" {
		return 0, errors.New("Внутренняя ошибка сервера")
	}
	return 1, nil
}

func (m *MockService) DeleteProduct(ctx context.Context, productId int) error {
	if productId == 2 {
		return errors.New("Внутренняя ошибка сервера")
	} else if productId == 3 {
		return repository.ErrNotFound
	}
	return nil
}
func (m *MockService) UpdateProduct(ctx context.Context, userId int, user *entities.ProductInfo) error {
	if user.Description == "интернал" {
		return errors.New("Внутренняя ошибка сервера")
	} else if user.Description == "не фаунд" {
		return repository.ErrNotFound
	}

	return nil
}

// Мок для кафки
type MockKafka struct{}

func (m *MockKafka) SendMessage(prdInfo entities.ProductInfo, action string) error {
	if prdInfo.ProductId == 5 {
		return errors.New("Внутренняя ошибка сервера")
	} else if prdInfo.Description == "kafka" {
		return errors.New("Внутренняя ошибка сервера")
	}
	return nil
}

func (m *MockKafka) Close() error {
	return nil
}

func TestHandler_AddNewProduct_Correct(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "корректненько",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	handler.addNewProduct(c)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	var response map[string]int
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response["productId"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_AddNewProduct_InorrectMissingPole(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		Category        string   `json:"category"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		Category:        "кино",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	handler.addNewProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key: 'ProductInfo.Description' Error:Field validation for 'Description' failed on the 'required' tag", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_AddNewProduct_InorrectValidation(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		Category:        " ",
		Description:     "корректненько",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	handler.addNewProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid category:   does not match regexp", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_AddNewProduct_CorrectButInternalErr(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "интернал",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	handler.addNewProduct(c)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Внутренняя ошибка сервера", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_UpdateProduct_Correct(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "корректненько",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "1"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	var response string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "OK", response)

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()

}

func TestHandler_UpdateProduct_CorrectButInternalErr(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "интернал",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "1"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Внутренняя ошибка сервера", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_UpdateProduct_CorrectButNotFound(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "не фаунд",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "1"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, repository.ErrNotFound.Error(), response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_UpdateProduct_IncorrectMissingProductId(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "не фаунд",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "productId parametr does not exist in path", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_UpdateProduct_IncorrectIncorrectProductId(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "не фаунд",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "kot"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "strconv.Atoi: parsing \"kot\": invalid syntax", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()

}

func TestHandler_UpdateProduct_IncorrectIncorrectValidationProductId(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "не фаунд",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "0"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid user id: can`t be less or equal 0", response["reason"])

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()

}

func TestHandler_UpdateProduct_IncorrectMissingField(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		Category string `json:"category"`
		Status   string `json:"status"`
	}

	product := ProductInfo{
		Category: "кино",
		Status:   "avaible",
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "1"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_DeleteProduct_Correct(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "1"},
	}

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

}

func TestHandler_DeleteProduct_CorrectButNotFound(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "3"},
	}

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

}

func TestHandler_DeleteProduct_CorrectButInternalError(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "2"},
	}

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

}

func TestHandler_DeleteProduct_IncorrectMissingParam(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

}

func TestHandler_DeleteProduct_IncorrectValidation(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "-1"},
	}

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

}

func TestHandler_DeleteProduct_IncorrectINcorrect(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "kot"},
	}

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

}

// Ошибка кафки
func TestHandler_DeleteProduct_IncorrectKafkaError(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "5"},
	}

	handler.deleteProduct(c)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestHandler_UpdateProduct_IncorrectKafkaError(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "корректненько",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	c.Params = gin.Params{
		gin.Param{Key: "productId", Value: "5"},
	}

	handler.updateProduct(c)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}

func TestHandler_AddNewUser_IncorrectKafkaError(t *testing.T) {

	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)
	type ProductInfo struct {
		ProductId       int      `json:"userId"`
		Category        string   `json:"category"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		ProductKeyWords []string `json:"productKeyWords"`
	}

	product := ProductInfo{
		ProductId:       1,
		Category:        "кино",
		Description:     "kafka",
		Status:          "avaible",
		ProductKeyWords: []string{"фильм"},
	}

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	reader := bytes.NewReader(jsonData)

	c.Request, err = http.NewRequest("POST", "/product", reader)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	handler.addNewProduct(c)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]int
	json.Unmarshal(w.Body.Bytes(), &response)

	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			t.Fatal()
		}
	}()
}
