package kafka

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/internal/service"
	"github.com/IBM/sarama"
)

type Consumer struct {
	Consumer sarama.Consumer
	topics   []string
	log      *slog.Logger
}

func NewConsumer(addrs []string, topics []string, log *slog.Logger) (*Consumer, error) {
	config := InitConfig()

	c, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer: c,
		topics:   topics,
		log:      log,
	}, nil
}

func (c *Consumer) Consume(handler service.KafkaHandler, ctx context.Context) error {
	fi := "kafka.Consumer.Consume"
	responseCh := make(chan *sarama.ConsumerMessage, 10)
	defer close(responseCh)
	//подписываемся на обновления топикоы
	for _, elem := range c.topics {
		go ConsumeTopic(c, elem, ctx, responseCh)
	}
	//обрабатываем ответы из топиков
	for {
		select {
		case msg, ok := <-responseCh:
			if ok {
				switch msg.Topic {
				case "user_updates":
					c.log.Info(fi + ": " + "Message about user_updates received from topic")
					if err := handler.AddUserData(ctx, msg); err != nil {
						c.log.Error("Error adding user data", "err", err)
						return err
					}
				case "product_updates":
					c.log.Info(fi + ": " + "Message about product_updates received from topic")
					if err := handler.AddProductData(ctx, msg); err != nil {
						c.log.Error("Error adding user data", "err", err)
						return err
					}
				}
			}
		case <-ctx.Done():
			c.log.Info("Closing all consumers by reason from server")
			return nil
		}
	}

}

func ConsumeTopic(
	c *Consumer, topic string,
	ctx context.Context, responseCh chan<- *sarama.ConsumerMessage,
) {
	pc, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		c.log.Error("Error consuming partition", "err", err)
	}
	defer pc.AsyncClose()

	for {
		select {
		case msg := <-pc.Messages():
			responseCh <- msg
			c.log.Info("From Message received from topic", topic, msg)
		case err := <-pc.Errors():
			c.log.Error("Error from consumer", err.Error(), err)
		case <-ctx.Done():
			c.log.Info("Closing consumer by reason from server, topic", topic, 0)
			return
		}

	}
}

func InitConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

func ConnectToKafka(loger *slog.Logger) *Consumer {
	fi := "main.connectToKafka"

	str := os.Getenv("KAFKA_ADDRS")
	tpc := os.Getenv("KAFKA_TOPIC")
	addrs := strings.Split(str, ",")
	topics := strings.Split(tpc, ",")

	c, err := NewConsumer(addrs, topics, loger)

	if err != nil {
		log.Panic(fi + ":" + err.Error())
	}

	return c
}
