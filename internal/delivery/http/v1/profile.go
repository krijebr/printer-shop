package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProfileHandlers struct {
	usecase usecase.User
}

func NewProfileHandlers(uc usecase.User) *ProfileHandlers {
	return &ProfileHandlers{usecase: uc}
}

func (h *UserHandlers) getProfile(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *UserHandlers) updateProfile(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
