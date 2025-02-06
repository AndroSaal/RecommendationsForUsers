package api

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitRoutes_Correct(t *testing.T) {
	handler := NewHandler(
		&MockService{},
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
		&MockKafka{},
	)

	router := handler.InitRoutes()
	assert.NotNil(t, router)
}
