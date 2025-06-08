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
var ErrProducerIsUsed = errors.New("this producer is used")
var ErrProductIsUsed = errors.New("this product is used")
var ErrCartIsEmpty = errors.New("cart is empty")
var ErrOrderNotFound = errors.New("order not found")
var ErrProductIsHidden = errors.New("product is hidden")
