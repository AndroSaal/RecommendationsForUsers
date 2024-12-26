package kafka

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/entities"
	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/product/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	Producer sarama.SyncProducer
	log      *slog.Logger
}

func NewProducer(brokerAdressses []string, log *slog.Logger) (*Producer, error) {

	producer, err := sarama.NewSyncProducer(brokerAdressses, InitConfig(brokerAdressses))
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: producer,
		log:      log,
	}, nil
}

func InitConfig(brokerAdressses []string) *sarama.Config {
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	return config
}

func (p *Producer) SendMessage(prdInfo entities.ProductInfo) error {
	topic := os.Getenv("KAFKA_TOPIC")

	if topic == "" {
		p.log.Error("KAFKA_TOPIC not set")
		return fmt.Errorf("environment KAFKA_TOPIC not set")
	}

	productKeyWords := make([]string, len(prdInfo.ProductKeyWords))

	for _, elem := range prdInfo.ProductKeyWords {
		productKeyWords = append(productKeyWords, fmt.Sprintf("%v", elem))
	}

	userMassage := myproto.ProductAction{
		ProductId:        int64(prdInfo.ProductId),
		ProductKeyWords: productKeyWords,
	}

	data, err := proto.Marshal(&userMassage)
	if err != nil {
		return err
	}

	partition, offset, err := p.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})

	if err != nil {
		p.log.Error(err.Error())
		return err
	} else {
		p.log.Info(fmt.Sprintf(
			"Message is sent to topic %s, partition %d, offset %d", topic, partition, offset,
		))
	}

	return nil
}
