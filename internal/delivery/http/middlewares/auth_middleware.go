package middlewares

import (
	"log"
	"net/http"

	"github.com/krijebr/printer-shop/internal/delivery/http/common"
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
			return c.NoContent(http.StatusUnauthorized)
		}
		user, err := a.u.Auth.ValidateToken(c.Request().Context(), getToken(authHeader))
		if err != nil {
			switch {
			case err == usecase.ErrInvalidToken:
				return c.NoContent(http.StatusUnauthorized)
			default:
				log.Println("Ошибка валидации токена", err)
				c.NoContent(http.StatusInternalServerError)
			}
		}
		c.Set(common.UserIdContextKey, user.Id)
		return next(c)
	}

}
func getToken(header string) string {
	return header[7:]
}
