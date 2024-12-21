package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndroSaal/RecommendationsForUsers/app/pkg/config"
	mylog "github.com/AndroSaal/RecommendationsForUsers/app/pkg/log"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/service"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/api"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/server"
)

func main() {
	//з агрузка переменных окружения
	env := config.MustLoadEnv()

	// логгер
	logger := mylog.NewLogger(env)

	// конфига
	cfg := config.MustLoadConfig()

	// коннект к бд (Маст)
	dbConn := repository.NewPostgresDB(cfg.DBConf)

	//TODO: Инициализация соединения к серверу почты

	// слой репозитория
	repository := repository.NewUserRepository(dbConn, logger)

	// слой сервиса
	service := service.NewUserService(repository)

	// транспортный слой
	handlers := api.NewHandler(service)

	// инициализация сервера
	srv, err := server.NewServer(cfg.SrvConf, handlers.InitRoutes(), logger)
	if err != nil {
		log.Fatal(err)
	}

	// обработка остановки по сигналу
	ctxSig, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM,
	)
	defer stop()

	// обработка остановки по таймауту
	ctxTim, cancel := context.WithTimeout(context.Background(), cfg.SrvConf.Timeout)
	defer cancel()

	// запуск сервера
	go func() {
		if err = srv.Run(); err != http.ErrServerClosed {
			fmt.Println(fmt.Errorf("error occured while running server: " + err.Error()))
		} else {
			return
		}
	}()

	// graceful shutdown
	for {
		select {
		case <-ctxTim.Done():
			logger.Info("Server Stopped by timout")
			srv.Stop(ctxTim)
			return
		case <-ctxSig.Done():
			logger.Info("Server stopped by system signall")
			srv.Stop(ctxSig)
			return
		}
	}
}
