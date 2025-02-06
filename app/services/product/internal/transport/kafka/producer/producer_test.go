package kafka

import (
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestKafka_ConnectToKafka_Correct(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	Producer := ConnectToKafka(logger)

	defer func() {
		if err := Producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	assert.NotNil(t, Producer)

}

func TestKafka_ConnectToKafka_Incorrect(t *testing.T) {

	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if err := os.Setenv("KAFKA_ADDRS", "localhost:9092"); err != nil {
		t.Error(err)
	}
	Producer := ConnectToKafka(logger)
	assert.Nil(t, Producer)

	defer func() {
		if err := Producer.Close(); err != nil {
			t.Error(err)
		}
	}()

}

func TestKafka_TestSendMessage_Correct(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := os.Setenv("KAFKA_ADDRS", "kafka-test-product-user:9094"); err != nil {
		t.Error(err)
	}

	Producer := ConnectToKafka(logger)

	product := entities.ProductInfo{
		ProductId: 1,
		ProductKeyWords: []string{
			"Chocolate",
		},
	}

	defer func() {
		if err := Producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	err := Producer.SendMessage(product, "delete")
	assert.NoError(t, err)
}

func TestKafka_SendMessage_Incorrect(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := os.Setenv("KAFKA_TOPIC", ""); err != nil {
		t.Error(err)
	}

	Producer := ConnectToKafka(logger)

	product := entities.ProductInfo{
		ProductId: 1,
		ProductKeyWords: []string{
			"Chocolate",
		},
	}

	defer func() {
		if err := Producer.Close(); err != nil {
			t.Error(err)
		}
	}()

	err := Producer.SendMessage(product, "update")
	assert.Error(t, err)
}
