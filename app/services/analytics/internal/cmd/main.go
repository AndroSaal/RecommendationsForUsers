package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/repository"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/service"
	kafka "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/consumer"
	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/config"
	mylog "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/pkg/log"
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
	dbConn.DB.Close()

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
	kafkaConn := kafka.ConnectToKafka(logger)
	kafkaConn.Consume(service, ctxSig)
	defer kafkaConn.Consumer.Close()

	// обработка остановки по таймауту
	ctxTim, cancel := context.WithTimeout(context.Background(), cfg.SrvConf.Timeout)
	defer cancel()
	for {
		select {
		case <-ctxTim.Done():
			logger.Info("Stopped by timout")
			return
		case <-ctxSig.Done():
			logger.Info("Stopped by system signall")
			return
		}
	}
}
