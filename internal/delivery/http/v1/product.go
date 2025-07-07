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
		if c.QueryParam("producer_id") != "" || c.QueryParam("product_status") != "" {
			filter = new(entity.ProductFilter)
		}
		if c.QueryParam("producer_id") != "" {
			producerId, err := uuid.Parse(c.QueryParam("producer_id"))
			if err != nil {
				slog.Debug("invalid producer id", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrInvalidRequestCode,
					Message: ErrInvalidRequestMessage,
				})
			}
			filter.ProducerId = &producerId
		}
		if c.QueryParam("product_status") != "" {
			validate := validator.New()
			err := validate.Var(c.QueryParam("product_status"), "oneof=published hidden")
			if err != nil {
				slog.Debug("validation error", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			productStatus := entity.ProductStatus(c.QueryParam("product_status"))
			filter.Status = &productStatus
		}
		var (
			userRole entity.UserRole
			ok       bool
		)
		userRole, ok = c.Get(UserRoleContextKey).(entity.UserRole)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		if userRole != entity.UserRoleAdmin {
			if filter != nil {
				if filter.Status != nil {
					if *filter.Status == entity.ProductStatusHidden {
						return c.JSON(http.StatusForbidden, ErrResponse{
							Error:   ErrForbiddenCode,
							Message: ErrForbiddenMessage,
						})
					}
				}
			}
			if filter == nil {
				filter = new(entity.ProductFilter)
			}
			publishedStatus := entity.ProductStatusPublished
			filter.Status = &publishedStatus
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
				slog.Debug("producer not found", slog.Any("error", err))
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
			slog.Debug("invalid product id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		product, err := p.usecase.GetById(c.Request().Context(), productId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Debug("product not found", slog.Any("error", err))
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
		var (
			userRole entity.UserRole
			ok       bool
		)
		userRole, ok = c.Get(UserRoleContextKey).(entity.UserRole)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		if userRole != entity.UserRoleAdmin {
			if product.Status != entity.ProductStatusPublished {
				return c.JSON(http.StatusForbidden, ErrResponse{
					Error:   ErrForbiddenCode,
					Message: ErrForbiddenMessage,
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
			slog.Debug("invalid product id", slog.Any("error", err))
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
		if requestData.Name == "" && requestData.Price == 0 && requestData.ProducerId == uuid.Nil && requestData.Status == "" {
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
				slog.Debug("product not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			case errors.Is(err, usecase.ErrProducerNotFound):
				slog.Debug("producer not found", slog.Any("error", err))
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
			slog.Debug("invalid product id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		err = p.usecase.DeleteById(c.Request().Context(), productId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Debug("product not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			case errors.Is(err, usecase.ErrProductIsUsed):
				slog.Debug("product is used", slog.Any("error", err))
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
