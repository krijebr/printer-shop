package repo

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrProducerNotFound = errors.New("producer not found")
var ErrTokenNotFound = errors.New("token not found")
