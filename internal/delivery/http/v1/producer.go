package v1

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
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
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		log.Println("Получение всех производителей")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, producers)
	}
}

func (p *ProducerHandlers) createProducer() echo.HandlerFunc {
	type request struct {
		Name        string `json:"name" validate:"required,max=30,min=3"`
		Description string `json:"description" validate:"required,max=300,min=5"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			log.Println("Ошибка чтения тела запроса ", err)
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			log.Println("Невалидные данные ", err)
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}

		producer := entity.Producer{
			Name:        requestData.Name,
			Description: requestData.Description,
		}
		newProducer, err := p.usecase.Create(c.Request().Context(), producer)
		if err != nil {
			log.Println("Ошибка создания производителя", err)
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		log.Println("Производитель создан")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, newProducer)
	}
}

func (p *ProducerHandlers) getProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		producerId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		producer, err := p.usecase.GetById(c.Request().Context(), producerId)
		if err != nil {
			switch {
			case err == usecase.ErrProducerNotFound:
				log.Println("Производитель с таким id не найден", err)
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				log.Println("Ошибка получения производителя", err)
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		log.Println("Производитель получен")
		return c.JSON(http.StatusOK, producer)
	}
}
func (p *ProducerHandlers) updateProducerById() echo.HandlerFunc {
	type request struct {
		Name        string `json:"name,omitempty" validate:"omitempty,max=30,min=3"`
		Description string `json:"description,omitempty" validate:"omitempty,max=300,min=5"`
	}
	return func(c echo.Context) error {
		producerId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		var requestData request
		err = c.Bind(&requestData)
		if err != nil {
			log.Println("Ошибка чтения тела запроса ", err)
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			log.Println("Невалидные данные ", err)
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}
		producer := entity.Producer{
			Id:          producerId,
			Name:        requestData.Name,
			Description: requestData.Description,
		}
		updatedProducer, err := p.usecase.Update(c.Request().Context(), producer)
		if err != nil {
			switch {
			case err == usecase.ErrProducerNotFound:
				log.Println("Производителя с таким id не найдено", err)
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				log.Println("Ошибка обновления производителя", err)
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		log.Println("Производитель обновлен")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, updatedProducer)
	}
}
func (p *ProducerHandlers) deleteProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		producerId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		err = p.usecase.DeleteById(c.Request().Context(), producerId)
		if err != nil {
			switch {
			case err == usecase.ErrProducerNotFound:
				log.Println("Производителя с таким id не найдено", err)
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			case err == usecase.ErrProducerIsUsed:
				log.Println("Производитель используется", err)
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProducerIsUsedCode,
					Message: ErrProducerIsUsedMessage,
				})
			default:
				log.Println("Ошибка удаления производителя", err)
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		log.Println("Производитель удален")
		return c.NoContent(http.StatusOK)
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
