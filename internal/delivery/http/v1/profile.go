package v1

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/delivery/http/common"
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
		userId, ok := c.Get(common.UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		user, err := p.usecase.GetById(c.Request().Context(), userId)
		if err != nil {
			log.Println("Ошибка полученя профиля", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		log.Printf("Профиль пользователя с id %s получен", userId)
		return c.JSON(http.StatusOK, user)
	}
}
func (p *ProfileHandlers) updateProfile() echo.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name,omitempty" validate:"max=25,min=3"`
		LastName  string `json:"last_name,omitempty" validate:"max=25,min=3"`
		Password  string `json:"password,omitempty" validate:"max=25,min=8"`
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

		userId, ok := c.Get(common.UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		user := entity.User{
			Id:           userId,
			FirstName:    requestData.FirstName,
			LastName:     requestData.LastName,
			PasswordHash: requestData.Password,
		}
		updatedUser, err := p.usecase.Update(c.Request().Context(), user)
		if err != nil {
			log.Println("Ошибка обновления профиля")
			return c.NoContent(http.StatusInternalServerError)
		}
		log.Printf("Профиль пользователя с id %s обновлен", userId)
		return c.JSON(http.StatusOK, updatedUser)
	}
}
func RegisterProfileRoutes(u usecase.User, g *echo.Group) {
	a := NewProfileHandlers(u)
	g.GET("", a.getProfile())
	g.PUT("", a.updateProfile())
}
