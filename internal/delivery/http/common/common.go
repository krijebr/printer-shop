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
	ErrProductIsUsedCode           = 13
	ErrCartIsEmptyCode             = 14
	ErrProductNotExistCode         = 15
	ErrOrderNotExistCode           = 16
	ErrOrderCantBeUpdatedCode      = 17
	ErrUserIsUsedCode              = 18
	ErrUserIsBlockedCode           = 19
	ErrOrderCantBeDeletedCode      = 20

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
	ErrProducerIsUsedMessage          = "this producer is already used and can't be deleted"
	ErrProductIsUsedMessage           = "this product is already used and can't be deleted"
	ErrCartIsEmptyMessage             = "cart is empty"
	ErrProductNotExistMessage         = "product with this id doesn't exist"
	ErrOrderNotExistMessage           = "order with this id doesn't exist"
	ErrOrderCantBeUpdatedMessage      = "order can't be updated"
	ErrUserIsUsedMessage              = "this user can't be deleted"
	ErrUserIsBlockedMessage           = "user id blocked"
	ErrOrderCantBeDeletedMessage      = "order can't be deleted"

	UserIdContextKey   string = "userId"
	UserRoleContextKey string = "userRole"
)
