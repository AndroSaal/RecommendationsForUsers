package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/repository"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockService struct{}

func (m *MockService) AddProductData(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return nil
}

func (m *MockService) AddUserData(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return nil
}

func (m *MockService) GetRecommendations(ctx context.Context, userId int) ([]int, error) {
	if userId == 2 {
		return nil, repository.ErrNotFound
	} else if userId == 3 {
		return nil, errors.New("Ошибка сервера")
	} else if userId == 4 {
		return nil, nil
	}
	return []int{1, 2, 3}, nil
}

func TestHandler_GetRecommendations_Correct(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "1"},
	}

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	var response []int
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, []int{1, 2, 3}, response)
}

func TestHandler_GetRecommendations_INcorrectMIssingParam(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/1", nil)

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "userId parameter is empty in path", response["reason"])
}

func TestHandler_GetRecommendations_INcorrectIncorrectParam(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "Kot"},
	}

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "userId parameter incorrect in path", response["reason"])
}

func TestHandler_GetRecommendations_INcorrectValidationParam(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "0"},
	}

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "userId validation failed: invalid user id: can`t be less or equel 0", response["reason"])
}

func TestHandler_GetRecommendations_CorrectButNotFoundUser(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/2", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "2"},
	}

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, repository.ErrNotFound.Error(), response["reason"])
}

func TestHandler_GetRecommendations_CorrectButINternalServerError(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/3", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "3"},
	}

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Ошибка сервера", response["reason"])
}

func TestHandler_GetRecommendations_CorrectButNotFountRecom(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// /recommendation/:userId

	// формируем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/recommendation/4", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userId", Value: "4"},
	}

	handler.getUserRecommendations(c)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "recommendations for this user not found", response["reason"])
}
