package repository

import "errors"

var (
	ErrNotFound = errors.New("user not found")
)

// var (
// 	//ошибки из базы
// 	errUniqueConstraintEmail string = `pq: duplicate key value violates unique constraint \"users_email_key\"`
// )
