package common

type ErrResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

const (
	ErrInvalidTokenCode            = 1
	ErrInvalidRefreshTokenCode     = 2
	ErrResourceNotFoundCode        = 3
	ErrInternalErrorCode           = 4
	ErrUnauthorizedCode            = 5
	ErrForbiddenCode               = 6
	ErrInvalidRequestCode          = 7
	ErrValidationErrorCode         = 8
	ErrEmailAlreadyExistsCode      = 9
	ErrProducerNotExistCode        = 10
	ErrInvalidLoginCredentialsCode = 11
	ErrProducerIsUsedCode          = 12

	ErrInvalidTokenMessage            = "invalid token"
	ErrInvalidRefreshTokenMessage     = "invalid token"
	ErrResourceNotFoundMessage        = "resource not found"
	ErrInternalErrorMessage           = "internal error"
	ErrUnauthorizedMessage            = "unauthorized"
	ErrForbiddenMessage               = "forbidden"
	ErrInvalidRequestMessage          = "invalid request"
	ErrValidationErrorMessage         = "validation error"
	ErrEmailAlreadyExistsMessage      = "user with this email already exists"
	ErrProducerNotExistMessage        = "producer with this id doesn't exist"
	ErrInvalidLoginCredentialsMessage = "wrong email or password"
	ErrProducerIsUsedMessage          = "this producer is already used and can not be deleted"

	UserIdContextKey string = "userId"
)
