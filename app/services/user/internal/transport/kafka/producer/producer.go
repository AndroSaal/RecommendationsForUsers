package kafka

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/kafka/pb"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type Producer interface {
	SendMessage(usrInfo entities.UserInfo) error
	Close() error
}

type KafkaProducer struct {
	Producer sarama.SyncProducer
	log      *slog.Logger
}

func NewProducer(brokerAdressses []string, log *slog.Logger) (Producer, error) {
	fi := "transport.kafka.NewProducer"
	producer, err := sarama.NewSyncProducer(brokerAdressses, InitConfig(brokerAdressses))
	if err != nil {
		log.Debug("%s: Error with making new producer: %v", fi, err)
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		log:      log,
	}, nil
}

func InitConfig(brokerAdressses []string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	return config
}

func (p *KafkaProducer) SendMessage(usrInfo entities.UserInfo) error {
	topic := os.Getenv("KAFKA_TOPIC")

	if topic == "" {
		p.log.Error("KAFKA_TOPIC not set")
		return fmt.Errorf("environment KAFKA_TOPIC not set")
	}

	uinterests := make([]string, len(usrInfo.UserInterests))

	for _, elem := range usrInfo.UserInterests {
		uinterests = append(uinterests, fmt.Sprintf("%v", elem))
	}

	userMassage := myproto.UserUpdate{
		UserId:        int64(usrInfo.UsrId),
		UserInterests: uinterests,
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

func (p *KafkaProducer) Close() error {
	err := p.Producer.Close()
	if err != nil {
		p.log.Error("error closing KafkaProducer" + err.Error())
	}
	return err
}

func ConnectToKafka(loger *slog.Logger) Producer {
	fi := "main.connectToKafka"

	str := os.Getenv("KAFKA_ADDRS")
	addrs := strings.Split(str, ",")

	p, err := NewProducer(addrs, loger)

	if err != nil {
		log.Fatal(fi + ":" + err.Error())
	}

	return p
}
