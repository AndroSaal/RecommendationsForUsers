package service

import "github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"

type Service interface {
	UserCreator
	UserGetter
	UserUpdator
	CodeVerifactor
}

type UserCreator interface {
	CreateUser(user *entities.UserInfo) (int, error)
}

type UserGetter interface {
	GetUserById(id int) (*entities.UserInfo, error)
	GetUserByEmail(email string) (*entities.UserInfo, error)
}

type UserUpdator interface {
	UpdateUser(userId int, user *entities.UserInfo) error
}

type MailSender interface {
	SendMail(email string, code string) error
}

type CodeVerifactor interface {
	VerifyCode(userId int, code string) (bool, error)
}
