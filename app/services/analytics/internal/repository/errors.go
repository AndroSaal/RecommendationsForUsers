package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("prodcut with such email already exists")
	ErrNotFound      = errors.New("product not found")
)

// var (
// 	//ошибки из базы
// 	errUniqueConstraintEmail string = `pq: duplicate key value violates unique constraint \"users_email_key\"`
// )
