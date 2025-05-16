package v1

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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
			userStatus := entity.UserStatus(c.QueryParam("user_status"))
			userFilter.UserStatus = &userStatus
		}
		if c.QueryParam("user_role") != "" {
			userRole := entity.UserRole(c.QueryParam("user_role"))
			userFilter.UserRole = &userRole
		}
		users, err := u.usecase.GetAll(c.Request().Context(), userFilter)
		if err != nil {
			log.Println("Ошибка получения ", err)
			return c.String(http.StatusInternalServerError, "")
		}
		log.Println("Получение всех пользователей")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, users)
	}
}

func (u *UserHandlers) register() echo.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
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

		newUser := entity.User{
			FirstName:    requestData.FirstName,
			LastName:     requestData.LastName,
			Email:        requestData.Email,
			PasswordHash: requestData.Password,
		}

		user, err := u.usecase.Register(c.Request().Context(), newUser)
		if err != nil {
			switch {
			case err == usecase.ErrEmailAlreadyExists:
				log.Println("Пользователь с таким email уже существует", err)
				return c.String(http.StatusBadRequest, "")
			default:
				log.Println("Ошибка создания ползьзователя", err)
				return c.String(http.StatusInternalServerError, "")
			}
		}
		log.Println("Регистрация нового пользователя")
		return c.JSON(http.StatusOK, user)
	}
}

func (u *UserHandlers) getUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
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
	g.POST("", a.register())
	g.GET("/:id", a.getUserById())
	g.PUT("/:id", a.updateUserById())
	g.DELETE("/:id", a.deleteUserById())
}
