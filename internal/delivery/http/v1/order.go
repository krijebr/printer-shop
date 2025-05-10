package v1

import (
	"net/http"

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
		return c.String(http.StatusNotImplemented, "Not Implemented")
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
