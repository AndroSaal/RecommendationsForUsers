package service

import (
	"context"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
)

type Service interface {
	UserCreator
	UserGetter
	UserUpdator
	CodeVerifactor
}

type UserCreator interface {
	CreateUser(ctx context.Context, user *entities.UserInfo) (int, error)
}

type UserGetter interface {
	GetUserById(ctx context.Context, id int) (*entities.UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error)
}

type UserUpdator interface {
	UpdateUser(ctx context.Context, userId int, user *entities.UserInfo) error
}

type MailSender interface {
	SendMail(ctx context.Context, email string, code string) error
}

type CodeVerifactor interface {
	VerifyCode(ctx context.Context, userId int, code string) (bool, error)
}
