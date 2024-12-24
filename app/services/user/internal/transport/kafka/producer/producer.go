package kafka

import (
	"fmt"
	"log/slog"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/transport/kafka/pb"
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
	}, nil
}

func InitConfig(brokerAdressses []string) *sarama.Config {
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	return config
}

func (p *Producer) SendMessage(topic string, usrInfo entities.UserInfo) error {
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

	p.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})

	return nil
}
