package mylog

import (
	"log"
	"log/slog"
	"os"
)

// создание логгера
func MustNewLogger(env string) *slog.Logger {

	fi := "NewLogger"

	var logger *slog.Logger
	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatal(fi + ":" + "Wrong evironment: " + env)
	}

	return logger
}
