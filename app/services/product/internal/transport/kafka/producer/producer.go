package kafka

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

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
	fi := "transport.kafka.NewProducer"
	producer, err := sarama.NewSyncProducer(brokerAdressses, InitConfig())
	if err != nil {
		log.Debug("%s: Error adding new user: %v", fi, err)
		return nil, err
	}

	return &Producer{
		Producer: producer,
		log:      log,
	}, nil
}

func InitConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	return config
}

func (p *Producer) SendMessage(prdInfo entities.ProductInfo, action string) error {
	fi := "transport.kafka.Producer.SendMessage"

	topic := os.Getenv("KAFKA_TOPIC")

	if topic == "" {
		p.log.Error("%s: Error Getting topic: %v", fi, errors.New("env KAFKA_TOPIC not set"))
		return fmt.Errorf("environment KAFKA_TOPIC not set")
	}

	productKeyWords := make([]string, len(prdInfo.ProductKeyWords))

	for _, elem := range prdInfo.ProductKeyWords {
		productKeyWords = append(productKeyWords, fmt.Sprintf("%v", elem))
	}

	userMassage := myproto.ProductAction{
		ProductId:       int64(prdInfo.ProductId),
		ProductKeyWords: productKeyWords,
		Action:          action,
	}

	data, err := proto.Marshal(&userMassage)
	if err != nil {
		p.log.Error("%s: Error Marshal struct userMassage to protobuf: %v", fi, err)
		return err
	}

	partition, offset, err := p.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})

	if err != nil {
		p.log.Error("%s: Error sending userMassage to kafka: %v", fi, err)
		return err
	} else {
		p.log.Info(fmt.Sprintf(
			"Message is sent to topic %s, partition %d, offset %d", topic, partition, offset,
		))
	}

	return nil
}

func ConnectToKafka(loger *slog.Logger) *Producer {
	fi := "main.connectToKafka"

	str := os.Getenv("KAFKA_ADDRS")
	addrs := strings.Split(str, ",")

	p, err := NewProducer(addrs, loger)

	if err != nil {
		log.Fatal(fi + ":" + err.Error())
	}

	return p
}
