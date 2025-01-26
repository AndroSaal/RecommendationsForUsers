package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("user with such email already exists")
	ErrNotFound      = errors.New("user not found")
)