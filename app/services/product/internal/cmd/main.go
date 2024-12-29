package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/service"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/transport/api"
	kafka "github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/transport/kafka/producer"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/transport/server"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/pkg/config"
	mylog "github.com/AndroSaal/RecommendationsForUsers/app/services/product/pkg/log"
)

func main() {
	//загрузка переменных окружения
	env := config.LoadEnv()

	// логгер
	logger := mylog.MustNewLogger(env)

	// конфига
	cfg := config.MustLoadConfig()

	// коннект к бд (Маст)
	dbConn := repository.NewPostgresDB(cfg.DBConf)
	defer dbConn.DB.Close()

	// слой репозитория
	repository := repository.NewProductRepository(dbConn, logger)

	// слой сервиса
	service := service.NewProductService(repository, logger)

	//коннект к кафке
	kafkaConn := kafka.ConnectToKafka(logger)
	defer kafkaConn.Producer.Close()

	// транспортный слой
	handlers := api.NewHandler(service, logger, kafkaConn)

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
