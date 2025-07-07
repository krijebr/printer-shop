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
			slog.Error("producers receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("all producers received")
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

		producer := entity.Producer{
			Name:        requestData.Name,
			Description: requestData.Description,
		}
		newProducer, err := p.usecase.Create(c.Request().Context(), producer)
		if err != nil {
			slog.Error("producer creation error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("producer created")
		return c.JSON(http.StatusOK, newProducer)
	}
}

func (p *ProducerHandlers) getProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		producerId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Debug("invalid producer id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		producer, err := p.usecase.GetById(c.Request().Context(), producerId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProducerNotFound):
				slog.Debug("producer not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				slog.Error("producer receiving error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("producer received")
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
			slog.Debug("invalid producer id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		var requestData request
		err = c.Bind(&requestData)
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
		if requestData.Name == "" && requestData.Description == "" {
			slog.Debug("validation error")
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
			case errors.Is(err, usecase.ErrProducerNotFound):
				slog.Debug("producer not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				slog.Error("producer updating", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("producer updated")
		return c.JSON(http.StatusOK, updatedProducer)
	}
}

func (p *ProducerHandlers) deleteProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		producerId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Debug("invalid producer id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		err = p.usecase.DeleteById(c.Request().Context(), producerId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProducerNotFound):
				slog.Debug("producer not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			case errors.Is(err, usecase.ErrProducerIsUsed):
				slog.Debug("producer is used", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProducerIsUsedCode,
					Message: ErrProducerIsUsedMessage,
				})
			default:
				slog.Error("producer delete error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("producer deleted")
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
