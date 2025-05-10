package v1

import (
	"net/http"

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
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (p *ProductHandlers) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (p *ProductHandlers) getProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProductHandlers) updateProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProductHandlers) deleteProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
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
