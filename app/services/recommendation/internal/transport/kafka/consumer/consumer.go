package kafka

import (
	"context"
	"log/slog"

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

func (c *Consumer) Consume(handler service.KafkaHandler, ctx context.Context) {
	responseCh := make(chan *sarama.ConsumerMessage, 10)
	defer close(responseCh)
	//подписываемся на обновления первого топика
	go ConsumeTopic(c, c.topics[0], ctx, responseCh)
	//подписываемся на обновления второго топика
	go ConsumeTopic(c, c.topics[1], ctx, responseCh)
	//обрабатываем ответы из топиков
	go func() {
		for {
			select {
			case msg := <-responseCh:
				switch msg.Topic {
				case "user_updates":
					handler.AddUserData(msg)
				case "product_updates":
					handler.AddProductData(msg)
				}
			case <-ctx.Done():
				c.log.Info("Closing all consumers by reason from server")
				return
			}
		}
	}()

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
			c.log.Info("From Message received from topic", topic)
		case err := <-pc.Errors():
			c.log.Error("Error from consumer", err)
		case <-ctx.Done():
			c.log.Info("Closing consumer by reason from server, topic", topic)
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
