package v1

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
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
		userId, ok := c.Get(UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		user, err := p.usecase.GetById(c.Request().Context(), userId)
		if err != nil {
			slog.Error("profile receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("profile received")
		return c.JSON(http.StatusOK, user)
	}
}
func (p *ProfileHandlers) updateProfile() echo.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name,omitempty" validate:"omitempty,max=25,min=3"`
		LastName  string `json:"last_name,omitempty" validate:"omitempty,max=25,min=3"`
		Password  string `json:"password,omitempty" validate:"omitempty,max=60,min=8"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			slog.Error("invalid request", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			slog.Error("validation error", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}
		if requestData.FirstName == "" && requestData.LastName == "" && requestData.Password == "" {
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}
		userId, ok := c.Get(UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		user := entity.User{
			Id:           userId,
			FirstName:    requestData.FirstName,
			LastName:     requestData.LastName,
			PasswordHash: requestData.Password,
		}
		updatedUser, err := p.usecase.Update(c.Request().Context(), user)
		if err != nil {
			slog.Error("profile updating error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("profile updated")
		return c.JSON(http.StatusOK, updatedUser)
	}
}
func RegisterProfileRoutes(u usecase.User, g *echo.Group) {
	a := NewProfileHandlers(u)
	g.GET("", a.getProfile())
	g.PUT("", a.updateProfile())
}
