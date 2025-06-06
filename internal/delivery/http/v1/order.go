package v1

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type OrderHandlers struct {
	usecase usecase.Order
}

func NewOrderHandlers(u usecase.Order) *OrderHandlers {
	return &OrderHandlers{usecase: u}
}

func (o *OrderHandlers) getAllOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter *entity.OrderFilter
		filter = nil
		if c.QueryParam("user_id") != "" || c.QueryParam("order_status") != "" {
			filter = new(entity.OrderFilter)
		}
		if c.QueryParam("user_id") != "" {

			userId, err := uuid.Parse(c.QueryParam("user_id"))
			if err != nil {
				slog.Error("invalid user id", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			filter.UserId = &userId
		}
		if c.QueryParam("order_status") != "" {
			validate := validator.New()
			err := validate.Var(c.QueryParam("order_status"), "oneof=new in_progress done")
			if err != nil {
				slog.Error("validation error", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			orderStatus := entity.OrderStatus(c.QueryParam("order_status"))
			filter.Status = &orderStatus
		}
		orders, err := o.usecase.GetAll(c.Request().Context(), filter)
		if err != nil {
			slog.Error("orders receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("all orders received")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, orders)
	}
}

func (o *OrderHandlers) createOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (o *OrderHandlers) getOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (o *OrderHandlers) updateOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (o *OrderHandlers) deleteOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterOrderRoutes(u usecase.Order, g *echo.Group) {
	a := NewOrderHandlers(u)
	g.GET("", a.getAllOrders())
	g.POST("", a.createOrder())
	g.GET("/:id", a.getOrderById())
	g.PUT("/:id", a.updateOrderById())
	g.DELETE("/:id", a.deleteOrderById())
}
