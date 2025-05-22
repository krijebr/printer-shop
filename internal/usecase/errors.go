package usecase

import (
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")
var ErrEmailAlreadyExists = errors.New("user with this email alredy exists")
var ErrUserNotFound = errors.New("user not found")
var ErrProducerNotFound = errors.New("producer not found")
var ErrProductNotFound = errors.New("product not found")
var ErrWrongPassword = errors.New("wrong password")
var ErrInvalidToken = errors.New("invalid token")
