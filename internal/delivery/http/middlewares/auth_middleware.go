package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/krijebr/printer-shop/internal/config"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

type AuthMiddleware struct {
	u       *usecase.UseCases
	roles   config.RoleConf
	baseUrl string
}

func NewAuthMiddleware(u *usecase.UseCases, r *config.RoleConf, baseUrl string) *AuthMiddleware {
	return &AuthMiddleware{
		u:       u,
		roles:   *r,
		baseUrl: baseUrl,
	}
}

func (a *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userRole entity.UserRole
			user     *entity.User
			err      error
		)
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			user, err = a.u.Auth.ValidateToken(c.Request().Context(), getToken(authHeader))
			if err != nil {
				if errors.Is(err, usecase.ErrInvalidToken) {
					slog.Debug("invalid token", slog.Any("error", err))
					return c.JSON(http.StatusUnauthorized, ErrResponse{
						Error:   ErrInvalidTokenCode,
						Message: ErrInvalidTokenMessage,
					})
				}
				slog.Error("token validation error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
			if user.Status == entity.UserStatusBlocked {
				return c.JSON(http.StatusForbidden, ErrResponse{
					Error:   ErrUserIsBlockedCode,
					Message: ErrUserIsBlockedMessage,
				})
			}
			c.Set(UserIdContextKey, user.Id)
			userRole = user.Role
			c.Set(UserRoleContextKey, userRole)
		} else {
			userRole = entity.UserRoleGuest
			c.Set(UserRoleContextKey, userRole)
		}
		path := strings.TrimPrefix(c.Path(), a.baseUrl)
		if methods, inMapPaths := a.roles[path]; inMapPaths {
			if roles, inMapMethods := methods[c.Request().Method]; inMapMethods {
				if slices.Contains(roles, string(userRole)) {
					return next(c)
				}
				if userRole == entity.UserRoleGuest {
					slog.Debug("unauthorized", slog.Any("error", err))
					return c.JSON(http.StatusUnauthorized, ErrResponse{
						Error:   ErrUnauthorizedCode,
						Message: ErrUnauthorizedMessage,
					})
				}
			}
		}
		return c.JSON(http.StatusForbidden, ErrResponse{
			Error:   ErrForbiddenCode,
			Message: ErrForbiddenMessage,
		})
	}

}

func getToken(header string) string {
	return header[7:]
}
