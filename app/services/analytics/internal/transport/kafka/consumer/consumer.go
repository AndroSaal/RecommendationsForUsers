package kafka

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/service"
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

func (c *Consumer) Consume(handler service.KafkaHandler, ctx context.Context) {
	responseCh := make(chan *sarama.ConsumerMessage, 10)
	defer close(responseCh)
	//подписываемся на обновления первого топика
	go ConsumeTopic(c, c.topics[0], ctx, responseCh)
	//подписываемся на обновления второго топика
	go ConsumeTopic(c, c.topics[1], ctx, responseCh)
	//обрабатываем ответы из топиков
	for {
		select {
		case msg, ok := <-responseCh:
			if ok {
				switch msg.Topic {
				case "user_updates":
					c.log.Info("Message about user_updates received from topic")
					handler.AddUserData(ctx, msg)
				case "product_updates":
					c.log.Info("Message about product_updates received from topic")
					handler.AddProductData(ctx, msg)
				}
			}
		case <-ctx.Done():
			c.log.Info("Closing all consumers by reason from server")
			return
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
		log.Fatal(fi + ":" + err.Error())
	}

	return c
}
