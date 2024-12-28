package service

import "github.com/IBM/sarama"

type Service interface {
	KafkaHandler
}


type KafkaHandler interface {
	AddProductData(msg *sarama.ConsumerMessage) error
	AddUserData(msg *sarama.ConsumerMessage) error
}
