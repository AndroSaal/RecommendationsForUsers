package service

import "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"

type Service interface {
	UserCreator
	UserGetter
	UserUpdator
	MailSender
}

type UserCreator interface {
	CreateUser(user entities.UserInfo) (entities.UserId, error)
}

type UserGetter interface {
	GetUserById(id entities.UserId) (entities.UserInfo, error)
	GetUserByEmail(email entities.Email) (entities.UserInfo, error)
}

type UserUpdator interface {
	UpdateUser(user entities.UserInfo) error
}

type MailSender interface {
	SendMail(email *entities.Email) error
}

func NewService() Service {
	return NewUserService()
}
