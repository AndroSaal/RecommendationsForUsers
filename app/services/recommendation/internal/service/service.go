package service

import (
	"context"

	"github.com/IBM/sarama"
)

type Service interface {
	RecomGetter
	KafkaHandler
}

type RecomGetter interface {
	GetRecommendations(ctx context.Context, userId int) ([]int, error)
}

type KafkaHandler interface {
	AddProductData(ctx context.Context, msg *sarama.ConsumerMessage) error
	AddUserData(ctx context.Context, msg *sarama.ConsumerMessage) error
}
