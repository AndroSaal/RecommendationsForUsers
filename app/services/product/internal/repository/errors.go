package repository

import "errors"

var (
	ErrNotFound = errors.New("product not found")
)

// var (
// 	//ошибки из базы
// 	errUniqueConstraintEmail string = `pq: duplicate key value violates unique constraint \"users_email_key\"`
// )
