package service

import "github.com/IBM/sarama"

type Service interface {
	RecomGetter
	KafkaHandler
}

type RecomGetter interface {
	GetRecommendations(userId int) ([]int, error)
}

type KafkaHandler interface {
	AddProductData(msg *sarama.ConsumerMessage) error
	AddUserData(msg *sarama.ConsumerMessage) error
}
