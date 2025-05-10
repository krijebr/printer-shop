package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandlers struct {
	usecase usecase.Auth
}

func NewAuthHandlers(u usecase.Auth) *AuthHandlers {
	return &AuthHandlers{usecase: u}
}
func (a *AuthHandlers) login() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (a *AuthHandlers) refreshTokens() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterAuthRoutes(u usecase.Auth, g *echo.Group) {
	a := NewAuthHandlers(u)
	g.POST("", a.login())
	g.POST("/refresh-tokens", a.refreshTokens())
}
