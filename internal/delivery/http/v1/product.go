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

type ProductHandlers struct {
	usecase usecase.Product
}

func NewProductHandlers(u usecase.Product) *ProductHandlers {
	return &ProductHandlers{usecase: u}
}

func (p *ProductHandlers) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter *entity.ProductFilter
		filter = nil
		if c.QueryParam("producer_id") != "" {
			filter = new(entity.ProductFilter)
			producerId, err := uuid.Parse(c.QueryParam("producer_id"))
			if err != nil {
				slog.Error("invalid producer id", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			}
			filter.ProducerId = &producerId
		}
		products, err := p.usecase.GetAll(c.Request().Context(), filter)
		if err != nil {
			slog.Error("products receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("all products received")
		return c.JSON(http.StatusOK, products)
	}
}

func (p *ProductHandlers) createProduct() echo.HandlerFunc {
	type request struct {
		Name       string               `json:"name" validate:"required,max=100,min=3"`
		Price      float32              `json:"price" validate:"required"`
		ProducerId uuid.UUID            `json:"producer_id" validate:"required,uuid"`
		Status     entity.ProductStatus `json:"status" validate:"required,oneof=published hidden"`
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
		product := entity.Product{
			Name:  requestData.Name,
			Price: requestData.Price,
			Producer: &entity.Producer{
				Id: requestData.ProducerId,
			},
			Status: requestData.Status,
		}
		newProduct, err := p.usecase.Create(c.Request().Context(), product)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProducerNotFound):
				slog.Error("producer not found", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProducerNotExistCode,
					Message: ErrProducerNotExistMessage,
				})
			default:
				slog.Error("product creation error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product created")
		return c.JSON(http.StatusOK, newProduct)
	}
}

func (p *ProductHandlers) getProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		productId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid product id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		product, err := p.usecase.GetById(c.Request().Context(), productId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Error("product not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				slog.Error("product receiving error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product received")
		return c.JSON(http.StatusOK, product)
	}
}
func (p *ProductHandlers) updateProductById() echo.HandlerFunc {
	type request struct {
		Name       string               `json:"name" validate:"omitempty,max=100,min=3"`
		Price      float32              `json:"price" validate:"omitempty"`
		ProducerId uuid.UUID            `json:"producer_id" validate:"omitempty,uuid"`
		Status     entity.ProductStatus `json:"status" validate:"omitempty,oneof=published hidden"`
	}
	return func(c echo.Context) error {
		productId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid product id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		var requestData request
		err = c.Bind(&requestData)
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

		product := entity.Product{
			Id:    productId,
			Name:  requestData.Name,
			Price: requestData.Price,
			Producer: &entity.Producer{
				Id: requestData.ProducerId,
			},
			Status: requestData.Status,
		}
		updatedProduct, err := p.usecase.Update(c.Request().Context(), product)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Error("product not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			case errors.Is(err, usecase.ErrProducerNotFound):
				slog.Error("producer not found", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProducerNotExistCode,
					Message: ErrProducerNotExistMessage,
				})
			default:
				slog.Error("product updating error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product updated")
		return c.JSON(http.StatusOK, updatedProduct)
	}
}
func (p *ProductHandlers) deleteProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		productId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid product id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		err = p.usecase.DeleteById(c.Request().Context(), productId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Error("product not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			case errors.Is(err, usecase.ErrProductIsUsed):
				slog.Error("product is used", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProductIsUsedCode,
					Message: ErrProductIsUsedMessage,
				})
			default:
				slog.Error("product delete error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product deleted")
		return c.NoContent(http.StatusOK)
	}
}
func RegisterProductRoutes(u usecase.Product, g *echo.Group) {
	a := NewProductHandlers(u)
	g.GET("", a.getAllProducts())
	g.POST("", a.createProduct())
	g.GET("/:id", a.getProductById())
	g.PUT("/:id", a.updateProductById())
	g.DELETE("/:id", a.deleteProductById())
}
