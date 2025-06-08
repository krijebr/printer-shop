package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CartHandlers struct {
	usecase usecase.Cart
}

func NewCartHandlers(u usecase.Cart) *CartHandlers {
	return &CartHandlers{usecase: u}
}
func (h *CartHandlers) getAllProductsInCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, ok := c.Get(UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		productsInCart, err := h.usecase.GetAllProducts(c.Request().Context(), userId)
		if err != nil {
			slog.Error("cart receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("cart received")
		return c.JSON(http.StatusOK, productsInCart)
	}
}
func (h *CartHandlers) addProductToCart() echo.HandlerFunc {
	type request struct {
		ProductId uuid.UUID `json:"product_id" validate:"required,uuid"`
		Count     int       `json:"count" validate:"min=0"`
	}
	return func(c echo.Context) error {
		userId, ok := c.Get(UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
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
		err = h.usecase.AddProduct(c.Request().Context(), userId, requestData.ProductId, requestData.Count)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Error("product not found", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProductNotExistCode,
					Message: ErrProductNotExistMessage,
				})
			case errors.Is(err, usecase.ErrProductIsHidden):
				slog.Error("product is hidden", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProductNotExistCode,
					Message: ErrProductNotExistMessage,
				})
			default:
				slog.Error("adding product to cart error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product added to cart")
		return c.NoContent(http.StatusOK)
	}
}
func RegisterCartRoutes(u usecase.Cart, g *echo.Group) {
	a := NewCartHandlers(u)
	g.GET("", a.getAllProductsInCart())
	g.POST("", a.addProductToCart())
}
