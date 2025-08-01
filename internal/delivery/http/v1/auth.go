package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandlers struct {
	usecase usecase.Auth
}

func NewAuthHandlers(u usecase.Auth) *AuthHandlers {
	return &AuthHandlers{usecase: u}
}

func (a *AuthHandlers) register() echo.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name" validate:"required,max=25,min=3"`
		LastName  string `json:"last_name" validate:"required,max=25,min=3"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,max=60,min=8"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			slog.Debug("invalid request", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			slog.Debug("validation error", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}

		user := entity.User{
			FirstName:    requestData.FirstName,
			LastName:     requestData.LastName,
			Email:        requestData.Email,
			PasswordHash: requestData.Password,
		}

		newUser, err := a.usecase.Register(c.Request().Context(), user)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrEmailAlreadyExists):
				slog.Debug("user with this email already exists", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrEmailAlreadyExistsCode,
					Message: ErrEmailAlreadyExistsMessage,
				})
			default:
				slog.Error("user creation error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("new user registered")
		return c.JSON(http.StatusOK, newUser)
	}
}

func (a *AuthHandlers) login() echo.HandlerFunc {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			slog.Debug("invalid request", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			slog.Debug("validation error", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}
		token, refreshToken, err := a.usecase.Login(c.Request().Context(), requestData.Email, requestData.Password)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrWrongPassword):
				slog.Debug("wrong email or password", slog.Any("error", err))
				return c.JSON(http.StatusForbidden, ErrResponse{
					Error:   ErrInvalidLoginCredentialsCode,
					Message: ErrInvalidLoginCredentialsMessage,
				})
			case errors.Is(err, usecase.ErrUserIsBlocked):
				slog.Debug("user is blocked", slog.Any("error", err))
				return c.JSON(http.StatusForbidden, ErrResponse{
					Error:   ErrUserIsBlockedCode,
					Message: ErrUserIsBlockedMessage,
				})
			default:
				slog.Error("authentication error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		responseData := response{
			Token:        token,
			RefreshToken: refreshToken,
		}
		slog.Info("user logged in")
		return c.JSON(http.StatusOK, responseData)
	}
}

func (a *AuthHandlers) refreshTokens() echo.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			slog.Debug("invalid request", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			slog.Debug("validation error", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}

		token, refreshToken, err := a.usecase.RefreshToken(c.Request().Context(), requestData.RefreshToken)
		if err != nil {
			if errors.Is(err, usecase.ErrInvalidToken) {
				slog.Debug("invalid refresh token", slog.Any("error", err))
				return c.JSON(http.StatusUnauthorized, ErrResponse{
					Error:   ErrInvalidRefreshTokenCode,
					Message: ErrInvalidRefreshTokenMessage,
				})
			}
			slog.Error("tokens refreshing error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		responseData := response{
			Token:        token,
			RefreshToken: refreshToken,
		}
		slog.Info("tokens refreshed")
		return c.JSON(http.StatusOK, responseData)
	}
}

func RegisterAuthRoutes(u usecase.Auth, g *echo.Group) {
	a := NewAuthHandlers(u)

	g.POST("", a.login())
	g.POST("/register", a.register())
	g.POST("/refresh-tokens", a.refreshTokens())
}
