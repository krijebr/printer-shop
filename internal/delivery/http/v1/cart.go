package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CartHandlers struct {
	usecase usecase.Cart
}

func NewCartHandlers(uc usecase.Cart) *CartHandlers {
	return &CartHandlers{usecase: uc}
}
func (h *CartHandlers) getAllProductsInCart(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *CartHandlers) addProductToCart(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
