package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandlers struct {
	usecase usecase.User
}

func NewUserHandlers(uc usecase.User) *UserHandlers {
	return &UserHandlers{usecase: uc}
}

func (h *UserHandlers) allUsers(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *UserHandlers) createUser(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *UserHandlers) getUserById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *UserHandlers) updateUserById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *UserHandlers) deleteUserById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
