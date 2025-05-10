package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandlers struct {
	usecase usecase.Auth
}

func NewAuthHandlers(uc usecase.Auth) *AuthHandlers {
	return &AuthHandlers{usecase: uc}
}
func (h *AuthHandlers) login(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *AuthHandlers) refreshTokens(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
