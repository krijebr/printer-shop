package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
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
		var userFilter *entity.UserFilter
		userFilter = nil
		if c.QueryParam("user_status") != "" || c.QueryParam("user_role") != "" {
			userFilter = new(entity.UserFilter)
		}
		if c.QueryParam("user_status") != "" {
			validate := validator.New()
			err := validate.Var(c.QueryParam("user_status"), "oneof=active blocked")
			if err != nil {
				slog.Error("validation error", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			userStatus := entity.UserStatus(c.QueryParam("user_status"))
			userFilter.UserStatus = &userStatus
		}

		if c.QueryParam("user_role") != "" {
			validate := validator.New()
			err := validate.Var(c.QueryParam("user_role"), "oneof=customer admin")
			if err != nil {
				slog.Error("validation error", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			userRole := entity.UserRole(c.QueryParam("user_role"))
			userFilter.UserRole = &userRole
		}
		users, err := u.usecase.GetAll(c.Request().Context(), userFilter)
		if err != nil {
			slog.Error("users receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("all users received")
		return c.JSON(http.StatusOK, users)
	}
}

func (u *UserHandlers) getUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid user id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		user, err := u.usecase.GetById(c.Request().Context(), userId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrUserNotFound):
				slog.Error("user not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				slog.Error("user receiving error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product received")
		return c.JSON(http.StatusOK, user)
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
	g.GET("/:id", a.getUserById())
	g.PUT("/:id", a.updateUserById())
	g.DELETE("/:id", a.deleteUserById())
}
