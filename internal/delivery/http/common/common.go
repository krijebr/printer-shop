package common

type ErrResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

const (
	ErrInvalidTokenCode                   = 1
	ErrInvalidTokenMessage                = "invalid token"
	ErrInvalidRefreshTokenCode            = 2
	ErrInvalidRefreshTokenMessage         = "invalid token"
	ErrResourceNotFoundCode               = 3
	ErrResourceNotFoundMessage            = "resource not found"
	ErrInternalErrorCode                  = 4
	ErrInternalErrorMessage               = "internal error"
	ErrUnauthorizedCode                   = 5
	ErrUnauthorizedMessage                = "unauthorized"
	ErrForbiddenCode                      = 6
	ErrForbiddenMessage                   = "forbidden"
	ErrInvalidRequestCode                 = 7
	ErrInvalidRequestMessage              = "invalid request"
	ErrValidationErrorCode                = 8
	ErrValidationErrorMessage             = "validation error"
	ErrEmailAlreadyExistsCode             = 9
	ErrEmailAlreadyExistsMessage          = "user with this email already exists"
	ErrProducerNotExistCode               = 10
	ErrProducerNotExistMessage            = "producer with this id doesn't exist"
	ErrWrongEmailOrPasswordCode           = 11
	ErrWrongEmailOrPasswordMessage        = "wrong email or password"
	ErrProducerIsUsedCode                 = 12
	ErrProducerIsUsedMessage              = "this producer is already used and can not be deleted"
	UserIdContextKey               string = "userId"
)
