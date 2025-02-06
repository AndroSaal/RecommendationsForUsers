package kafka

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// тестирование Consumer кафки с моками вместо топик-хэндлеров
type MockTopicHandler struct{}

func (m *MockTopicHandler) AddProductData(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var product myproto.ProductAction

	err := proto.Unmarshal(msg.Value, &product)
	if err != nil {
		return err
	}

	if product.ProductId == 400 {
		return errors.New("some error")
	}
	return nil
}

func (m *MockTopicHandler) AddUserData(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var user myproto.UserUpdate

	err := proto.Unmarshal(msg.Value, &user)
	if err != nil {
		return err
	}

	if user.UserId == 400 {
		return errors.New("some error")
	}
	return nil
}

func NewProducer(t *testing.T, addr []string, topics []string) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(addr, InitConfigProducer(addr))
	if err != nil {
		t.Error(err)
		return nil
	}

	return producer
}

func InitConfigProducer(brokerAdressses []string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	return config
}

func TestKafka_ConnectToKafka_Correct(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	consumer := ConnectToKafka(logger)

	assert.NotNil(t, consumer)
	defer func() {
		err := consumer.Consumer.Close()
		assert.NoError(t, err)
	}()

}

func TestKafka_ConnectToKafka_Incorrect(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	err1 := os.Setenv("KAFKA_ADDRS", "")
	err2 := os.Setenv("KAFKA_TOPIC", "")
	if err1 != nil || err2 != nil {
		t.Error(err1.Error(), err2.Error())
	}
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()

	consumer := ConnectToKafka(logger)

	assert.Nil(t, consumer)
	defer func() {
		err := consumer.Consumer.Close()
		assert.NoError(t, err)
	}()

}

func TestKafka_Consume_Correct(t *testing.T) {

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	//создаем продюсер, которым отправим сообщения в топики
	producer := NewProducer(t, []string{"kafka-test-analytics:9092"}, []string{"user_updates", "product_updates"})
	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()
	//создаем тестируемый консьюмер
	consumer, err := NewConsumer(
		[]string{"kafka-test-analytics:9092"}, []string{"user_updates", "product_updates"}, logger,
	)
	defer func() {
		if err := consumer.Consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	//если ошибка - тест валится
	assert.NoError(t, err)

	tim, _ := context.WithTimeout(context.Background(), 5*time.Second)
	go func() {
		err := consumer.Consume(&MockTopicHandler{}, tim)
		assert.NoError(t, err)
	}()

	//создаем и сериализуем сообщение для отправки
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        1,
		UserInterests: []string{"test"},
	}

	dataUser, errUser := proto.Marshal(&user)
	assert.NoError(t, errUser)

	//отправляем сообщение
	producer.SendMessage(&sarama.ProducerMessage{
		Topic: "user_updates",
		Value: sarama.StringEncoder(dataUser),
	})

	var product myproto.ProductAction = myproto.ProductAction{
		ProductId: 1,
		Action:    "test",
	}

	dataProduct, errProduct := proto.Marshal(&product)
	assert.NoError(t, errProduct)

	producer.SendMessage(&sarama.ProducerMessage{
		Topic: "product_updates",
		Value: sarama.StringEncoder(dataProduct),
	})

	time.Sleep(5 * time.Second)

}

func TestKafka_Consume_IncorrectProduct(t *testing.T) {

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	//создаем продюсер, которым отправим сообщения в топики
	producer := NewProducer(t, []string{"kafka-test-analytics:9092"}, []string{"product_updates"})
	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()
	//создаем тестируемый консьюмер
	consumer, err := NewConsumer(
		[]string{"kafka-test-analytics:9092"}, []string{"product_updates"}, logger,
	)
	defer func() {
		if err := consumer.Consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	//если ошибка - тест валится
	assert.NoError(t, err)

	tim, _ := context.WithTimeout(context.Background(), 2*time.Second)
	go func() {
		err := consumer.Consume(&MockTopicHandler{}, tim)
		assert.Error(t, err)
	}()

	//создаем и сериализуем сообщение для отправки
	var product myproto.ProductAction = myproto.ProductAction{
		ProductId: 400,
		Action:    "test",
	}

	//отправляем сообщение

	dataProduct, errProduct := proto.Marshal(&product)
	assert.NoError(t, errProduct)

	producer.SendMessage(&sarama.ProducerMessage{
		Topic: "product_updates",
		Value: sarama.StringEncoder(dataProduct),
	})

	time.Sleep(2 * time.Second)

}

func TestKafka_Consume_IncorrectUser(t *testing.T) {

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	//создаем продюсер, которым отправим сообщения в топики
	producer := NewProducer(t, []string{"kafka-test-analytics:9092"}, []string{"user_updates"})
	defer func() {
		if err := producer.Close(); err != nil {
			t.Error(err)
		}
	}()
	//создаем тестируемый консьюмер
	consumer, err := NewConsumer(
		[]string{"kafka-test-analytics:9092"}, []string{"user_updates"}, logger,
	)
	defer func() {
		if err := consumer.Consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	//если ошибка - тест валится
	assert.NoError(t, err)

	tim, _ := context.WithTimeout(context.Background(), 2*time.Second)
	go func() {
		err := consumer.Consume(&MockTopicHandler{}, tim)
		assert.Error(t, err)
	}()

	//создаем и сериализуем сообщение для отправки
	var user myproto.UserUpdate = myproto.UserUpdate{
		UserId:        400,
		UserInterests: []string{"test"},
	}

	dataUser, errUser := proto.Marshal(&user)
	assert.NoError(t, errUser)

	//отправляем сообщение
	producer.SendMessage(&sarama.ProducerMessage{
		Topic: "user_updates",
		Value: sarama.StringEncoder(dataUser),
	})

	time.Sleep(2 * time.Second)

}
