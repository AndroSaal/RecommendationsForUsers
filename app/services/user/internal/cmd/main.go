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

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/service"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/api"
	kafka "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/kafka/producer"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/server"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/config"
	mylog "github.com/AndroSaal/RecommendationsForUsers/app/services/user/pkg/log"
)

func main() {
	//з агрузка переменных окружения
	env := config.MustLoadEnv()

	// логгер
	logger := mylog.MustNewLogger(env)

	// конфига
	cfg := config.MustLoadConfig()

	// коннект к бд (Маст)
	dbConn := repository.NewPostgresDB(cfg.DBConf)

	//Инициализация соединения к серверу почты
	mail := service.NewMailSender(cfg.MailConf, logger)

	// слой репозитория
	repository := repository.NewUserRepository(dbConn, logger)

	// слой сервиса
	service := service.NewUserService(mail, repository, logger)

	//коннект к кафке
	kafkaConn := connectToKafka(logger)

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

func connectToKafka(loger *slog.Logger) *kafka.Producer {
	fi := "main.connectToKafka"

	str := os.Getenv("KAFKA_ADDRS")
	addrs := strings.Split(str, ",")

	p, err := kafka.NewProducer(addrs, loger)

	if err != nil {
		log.Fatal(fi + ":" + err.Error())
	}

	return p
}
