//основные сущности и их валидация

package entities

import (
	"errors"
	"time"

	myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

type UserFullUpdate struct {
	User      *myproto.UserUpdate `json:"user"`
	Timestamp time.Time           `json:"timestamp"`
}

type ProductFullUpdate struct {
	Product   *myproto.ProductAction `json:"product"`
	Timestamp time.Time              `json:"timestamp"`
}

func ValidateUserId(prId int) error {

	if prId < 0 {
		return errors.New("invalid user id: can`t be less 0")
	}

	return nil
}
