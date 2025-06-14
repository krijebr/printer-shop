package middlewares

import (
	"log/slog"
	"net/http"

	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	u *usecase.UseCases
}

func NewAuthMiddleware(u *usecase.UseCases) *AuthMiddleware {
	return &AuthMiddleware{
		u: u,
	}
}

func (a *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, ErrResponse{
				Error:   ErrUnauthorizedCode,
				Message: ErrUnauthorizedMessage,
			})
		}
		user, err := a.u.Auth.ValidateToken(c.Request().Context(), getToken(authHeader))
		if err != nil {
			switch {
			case err == usecase.ErrInvalidToken:
				return c.JSON(http.StatusUnauthorized, ErrResponse{
					Error:   ErrInvalidTokenCode,
					Message: ErrInvalidTokenMessage,
				})
			default:
				slog.Error("token validation error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		if user.Status == entity.UserStatusBlocked {
			return c.JSON(http.StatusForbidden, ErrResponse{
				Error:   ErrForbiddenCode,
				Message: ErrForbiddenMessage,
			})
		}
		c.Set(UserIdContextKey, user.Id)
		c.Set(UserRoleContextKey, user.Role)
		return next(c)
	}

}
func getToken(header string) string {
	return header[7:]
}
