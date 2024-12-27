package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/service"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/api"
	kafka "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/kafka/consumer"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/transport/server"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	mylog "github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/log"
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

	// слой репозитория
	repository := repository.NewProductRepository(dbConn, logger)

	// слой сервиса
	service := service.NewRecommendationService(repository, logger)

	// обработка остановки по сигналу
	ctxSig, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM,
	)
	defer stop()

	//коннект к кафке
	kafkaConn := connectToKafka(logger)
	kafkaConn.Consume(service, ctxSig)

	// транспортный слой
	handlers := api.NewHandler(service, logger)

	// инициализация сервера
	srv, err := server.NewServer(cfg.SrvConf, handlers.InitRoutes(), logger)
	if err != nil {
		log.Fatal(err)
	}

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

func connectToKafka(loger *slog.Logger) *kafka.Consumer {
	fi := "main.connectToKafka"

	str := os.Getenv("KAFKA_ADDRS")
	tpc := os.Getenv("KAFKA_TOPIC")
	addrs := strings.Split(str, ",")
	topics := strings.Split(tpc, ",")

	c, err := kafka.NewConsumer(addrs, topics, loger)

	if err != nil {
		log.Fatal(fi + ":" + err.Error())
	}

	return c
}
