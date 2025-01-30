package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestMailSender_NewMailSender_Correct(t *testing.T) {
	config := config.ServerMailConf{
		Login:    "test@example.com",
		Password: "testpassword",
		Host:     "smtp.example.com",
		Port:     "587",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mailSender := NewMailSender(config, logger)
	assert.NotNil(t, mailSender)
}

func TestMailSender_SendMail_Correct(t *testing.T) {

	config := config.ServerMailConf{
		Login:    "salnikandro@gmail.com",
		Password: "hpvknzrsdkkwwbeq",
		Host:     "smtp.gmail.com",
		Port:     "465",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mailSender := NewMailSender(config, logger)

	// Пустой email и тело письма
	err := mailSender.SendMail(context.Background(), "89251409398@mail.ru", "Test-Test")

	assert.NoError(t, err)
}

func TestMailSender_SendMail_IncorrectEmptyEmail(t *testing.T) {

	config := config.ServerMailConf{
		Login:    "salnikandro@gmail.com",
		Password: "hpvknzrsdkkwwbeq",
		Host:     "smtp.gmail.com",
		Port:     "465",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mailSender := NewMailSender(config, logger)

	// Пустой email и тело письма
	err := mailSender.SendMail(context.Background(), "", "Test-Test")

	assert.Error(t, err)
}

func TestMailSender_SendMail_IncorrectEmptuPort(t *testing.T) {

	config := config.ServerMailConf{
		Login:    "salnikandro@gmail.com",
		Password: "hpvknzrsdkkwwbeq",
		Host:     "smtp.gmail.com",
		Port:     "",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mailSender := NewMailSender(config, logger)

	// Пустой email и тело письма
	err := mailSender.SendMail(context.Background(), "89251409398@mail.ru", "Test-Test")

	assert.Error(t, err)
}

func TestMailSender_SendMail_IncorrectEmptuHost(t *testing.T) {

	config := config.ServerMailConf{
		Login:    "salnikandro@gmail.com",
		Password: "hpvknzrsdkkwwbeq",
		Host:     "test.gmail.com",
		Port:     "465",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mailSender := NewMailSender(config, logger)

	// Пустой email и тело письма
	err := mailSender.SendMail(context.Background(), "89251409398@mail.ru", "Test-Test")

	assert.Error(t, err)
}

func TestMailSender_SendMail_IncorrectLogin(t *testing.T) {

	config := config.ServerMailConf{
		Login:    "tests@gmail.com",
		Password: "hpvknzrsdkkwwbeq",
		Host:     "smtp.gmail.com",
		Port:     "465",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mailSender := NewMailSender(config, logger)

	// Пустой email и тело письма
	err := mailSender.SendMail(context.Background(), "89251409398@mail.ru", "Test-Test")

	assert.Error(t, err)
}
