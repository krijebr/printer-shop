package middlewares

import (
	"net/http"
	"strings"

	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

type RoleCheckMiddleware struct {
	roles map[string]map[string][]string
}

func NewRoleCheckMiddleware(r *map[string]map[string][]string) *RoleCheckMiddleware {
	return &RoleCheckMiddleware{
		roles: *r,
	}
}

func (r *RoleCheckMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, ok := c.Get(UserRoleContextKey).(entity.UserRole)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		path := strings.Trim(c.Path(), "/api/v1")
		if methods, inMap := r.roles[path]; inMap {
			if roles, inMap := methods[c.Request().Method]; inMap {
				if slices.Contains(roles, string(userRole)) {
					return next(c)
				}
			}
		}
		return c.JSON(http.StatusForbidden, ErrResponse{
			Error:   ErrForbiddenCode,
			Message: ErrForbiddenMessage,
		})
	}

}
