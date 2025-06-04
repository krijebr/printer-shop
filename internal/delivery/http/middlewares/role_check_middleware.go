package middlewares

import (
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type RoleCheckMiddleware struct {
	u *usecase.UseCases
}

func NewRoleCheckMiddleware(u *usecase.UseCases) *RoleCheckMiddleware {
	return &RoleCheckMiddleware{
		u: u,
	}
}

func (a *RoleCheckMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		return next(c)
	}

}
