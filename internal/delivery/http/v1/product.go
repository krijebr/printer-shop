package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProductHandlers struct {
	usecase usecase.Product
}

func NewProductHandlers(uc usecase.Product) *ProductHandlers {
	return &ProductHandlers{usecase: uc}
}

func (h *ProductHandlers) getAllProducts(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *ProductHandlers) createProduct(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *ProductHandlers) getProductById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *ProductHandlers) updateProductById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *ProductHandlers) deleteProductById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
