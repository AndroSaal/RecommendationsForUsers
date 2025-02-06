package repository

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/config"
	"github.com/stretchr/testify/assert"
)

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
	cfg := loadConf()

	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	user := myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"test"},
	}

	timestamp, err := dbConn.AddUserUpdate(context.Background(), &user)

	assert.NoError(t, err)
	assert.NotEmpty(t, timestamp)
}

func TestPostgreDB_AddUserUpdate_Incorrect(t *testing.T) {
	cfg := loadConf()

	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	user := myproto.UserUpdate{
		UserId:        0,
		UserInterests: []string{"test"},
	}

	timestamp, err := dbConn.AddUserUpdate(context.Background(), &user)

	assert.Error(t, err)
	assert.Empty(t, timestamp)
}

func TestPostgreDB_AddProductUpdate_Correct(t *testing.T) {
	cfg := loadConf()

	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	product := myproto.ProductAction{
		ProductId:       1,
		ProductKeyWords: []string{"test"},
	}

	timestamp, err := dbConn.AddProductUpdate(context.Background(), &product)

	assert.NoError(t, err)
	assert.NotEmpty(t, timestamp)
}

func TestPostgreDB_AddProductUpdate_Incorrect(t *testing.T) {
	cfg := loadConf()

	dbConn := NewPostgresDB(cfg, slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	product := myproto.ProductAction{
		ProductId:       0,
		ProductKeyWords: []string{"test"},
	}

	timestamp, err := dbConn.AddProductUpdate(context.Background(), &product)

	assert.Error(t, err)
	assert.Empty(t, timestamp)
}
