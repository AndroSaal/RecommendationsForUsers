package mylog

import (
	"log/slog"
	"os"
)

// создание логгера
func MustLoadNewLogger(env string) *slog.Logger {

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
		panic(fi + ":" + "Wrong evironment " + env)
	}

	return logger
}
