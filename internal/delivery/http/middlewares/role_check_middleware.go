package middlewares

import (
	"net/http"
	"strings"

	"github.com/krijebr/printer-shop/internal/config"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

type RoleCheckMiddleware struct {
	roles   config.RoleConf
	baseUrl string
}

func NewRoleCheckMiddleware(r *config.RoleConf, baseUrl string) *RoleCheckMiddleware {
	return &RoleCheckMiddleware{
		roles:   *r,
		baseUrl: baseUrl,
	}
}

func (r *RoleCheckMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userRole entity.UserRole
			ok       bool
		)
		if val := c.Get(UserRoleContextKey); val != nil {
			userRole, ok = c.Get(UserRoleContextKey).(entity.UserRole)
			if !ok {
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		} else {
			userRole = entity.UserRoleGuest
		}
		path := strings.TrimPrefix(c.Path(), r.baseUrl)
		if methods, inMapPaths := r.roles[path]; inMapPaths {
			if roles, inMapMethods := methods[c.Request().Method]; inMapMethods {
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
