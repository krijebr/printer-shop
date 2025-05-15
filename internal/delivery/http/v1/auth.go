package v1

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandlers struct {
	usecase usecase.Auth
}

func NewAuthHandlers(u usecase.Auth) *AuthHandlers {
	return &AuthHandlers{usecase: u}
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
			log.Println("Ошибка чтения тела запроса ", err)
			return c.String(http.StatusBadRequest, "")
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			log.Println("Не валидные данные ", err)
			return c.String(http.StatusBadRequest, "")
		}
		token, refreshToken, err := a.usecase.Login(c.Request().Context(), requestData.Email, requestData.Password)
		if err != nil {
			log.Println("Ошибка авторизации", err)
			return c.String(http.StatusInternalServerError, "")
		}
		responseData := response{
			Token:        token,
			RefreshToken: refreshToken,
		}
		return c.JSON(http.StatusOK, responseData)
	}
}
func (a *AuthHandlers) refreshTokens() echo.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token" validate:"jwt"`
	}
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			log.Println("Ошибка чтения тела запроса ", err)
			return c.String(http.StatusBadRequest, "")
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			log.Println("Не валидные данные ", err)
			return c.String(http.StatusBadRequest, "")
		}

		token, refreshToken, err := a.usecase.RefreshToken(c.Request().Context(), requestData.RefreshToken)
		if err != nil {
			log.Println("Ошибка обновления токенов", err)
			return c.String(http.StatusInternalServerError, "")
		}
		responseData := response{
			Token:        token,
			RefreshToken: refreshToken,
		}
		return c.JSON(http.StatusOK, responseData)
	}
}
func RegisterAuthRoutes(u usecase.Auth, g *echo.Group) {
	a := NewAuthHandlers(u)
	g.POST("", a.login())
	g.POST("/refresh-tokens", a.refreshTokens())
}
