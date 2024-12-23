package kafka

import "github.com/IBM/sarama"

type Producer struct {
	Producer sarama.SyncProducer
}

func NewProducer(brokerAdressses []string) (*Producer, error) {
	config := InitConfig(brokerAdressses)
	producer, err := sarama.NewSyncProducer(brokerAdressses, config)
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

func (p *Producer) SendMessage(topic string, key, value []byte) error {
	p.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
}
