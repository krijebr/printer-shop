package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandlers struct {
	usecase usecase.User
}

func NewUserHandlers(u usecase.User) *UserHandlers {
	return &UserHandlers{usecase: u}
}

func (u *UserHandlers) allUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (u *UserHandlers) createUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (u *UserHandlers) getUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (u *UserHandlers) updateUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (u *UserHandlers) deleteUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterUserRoutes(u usecase.User, g *echo.Group) {
	a := NewUserHandlers(u)
	g.GET("", a.allUsers())
	g.POST("", a.createUser())
	g.GET("/:id", a.getUserById())
	g.PUT("/:id", a.updateUserById())
	g.DELETE("/:id", a.deleteUserById())
}
