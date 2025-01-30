package api

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHandler_InitRoutes_Correct(t *testing.T) {
	// Создаем наш хэндлер (Собственно транспортный слой)

	handler := NewHandler(
		&MockService{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	// инициализируем маршруты
	router := handler.InitRoutes()

	// проверяем что маршруты инициализировались
	assert.NotNil(t, router)
}
