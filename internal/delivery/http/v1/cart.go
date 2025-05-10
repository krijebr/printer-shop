package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CartHandlers struct {
	usecase usecase.Cart
}

func NewCartHandlers(u usecase.Cart) *CartHandlers {
	return &CartHandlers{usecase: u}
}
func (c *CartHandlers) getAllProductsInCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (c *CartHandlers) addProductToCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterCartRoutes(u usecase.Cart, g *echo.Group) {
	a := NewCartHandlers(u)
	g.GET("", a.getAllProductsInCart())
	g.POST("", a.addProductToCart())
}
