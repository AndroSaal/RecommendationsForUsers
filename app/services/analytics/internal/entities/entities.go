//основные сущности и их валидация

package entities

import (
	"errors"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func ValidateUserId(prId int) error {

	if prId < 0 {
		return errors.New("invalid user id: can`t be less 0")
	}

	return nil
}
