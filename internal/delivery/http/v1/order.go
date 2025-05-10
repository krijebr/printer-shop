package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type OrderHandlers struct {
	usecase usecase.Order
}

func NewOrderHandlers(uc usecase.Order) *OrderHandlers {
	return &OrderHandlers{usecase: uc}
}

func (h *OrderHandlers) getAllOrders(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *OrderHandlers) createOrder(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *OrderHandlers) getOrderById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *OrderHandlers) updateOrderById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *OrderHandlers) deleteOrderById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
