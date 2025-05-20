package v1

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProducerHandlers struct {
	usecase usecase.Producer
}

func NewProducerHandlers(u usecase.Producer) *ProducerHandlers {
	return &ProducerHandlers{usecase: u}
}

func (p *ProducerHandlers) getAllProducers() echo.HandlerFunc {
	return func(c echo.Context) error {
		producers, err := p.usecase.GetAll(c.Request().Context())
		if err != nil {
			log.Println("Ошибка получения производителей", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		log.Println("Получение всех производителей")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, producers)
	}
}

func (p *ProducerHandlers) createProducer() echo.HandlerFunc {
	type request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			log.Println("Ошибка чтения тела запроса ", err)
			return c.NoContent(http.StatusBadRequest)
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			log.Println("Не валидные данные ", err)
			return c.String(http.StatusBadRequest, "")
		}

		producer := entity.Producer{
			Name:        requestData.Name,
			Description: requestData.Description,
		}
		newProducer, err := p.usecase.Create(c.Request().Context(), producer)
		if err != nil {
			log.Println("Ошибка создания производителя", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, newProducer)
	}
}

func (p *ProducerHandlers) getProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		producerId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.NoContent(http.StatusBadRequest)
		}
		producer, err := p.usecase.GetById(c.Request().Context(), producerId)
		if err != nil {
			switch {
			case err == usecase.ErrProducerNotFound:
				log.Println("Производителя с таким id не найдено", err)
				return c.NoContent(http.StatusBadRequest)
			default:
				log.Println("Ошибка получения производителя", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		return c.JSON(http.StatusOK, producer)
	}
}
func (p *ProducerHandlers) updateProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProducerHandlers) deleteProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterProducerRoutes(u usecase.Producer, g *echo.Group) {
	a := NewProducerHandlers(u)
	g.GET("", a.getAllProducers())
	g.POST("", a.createProducer())
	g.GET("/:id", a.getProducerById())
	g.PUT("/:id", a.updateProducerById())
	g.DELETE("/:id", a.deleteProducerById())
}
