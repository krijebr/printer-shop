package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProfileHandlers struct {
	usecase usecase.User
}

func NewProfileHandlers(u usecase.User) *ProfileHandlers {
	return &ProfileHandlers{usecase: u}
}

func (p *ProfileHandlers) getProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProfileHandlers) updateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterProfileRoutes(u usecase.User, g *echo.Group) {
	a := NewProfileHandlers(u)
	g.GET("", a.getProfile())
	g.PUT("", a.updateProfile())
}
