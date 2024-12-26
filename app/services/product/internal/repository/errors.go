package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("user with such email already exists")
	ErrNotFound      = errors.New("user not found")
)

// var (
// 	//ошибки из базы
// 	errUniqueConstraintEmail string = `pq: duplicate key value violates unique constraint \"users_email_key\"`
// )
