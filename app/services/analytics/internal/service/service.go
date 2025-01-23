package service

import (
	"context"

	"github.com/IBM/sarama"
)

type Service interface {
	KafkaHandler
}

type KafkaHandler interface {
	AddProductData(ctx context.Context, msg *sarama.ConsumerMessage) error
	AddUserData(ctx context.Context, msg *sarama.ConsumerMessage) error
}
