package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProducerHandlers struct {
	usecase usecase.Producer
}

func NewProducerHandlers(u usecase.Producer) *ProducerHandlers {
	return &ProducerHandlers{usecase: u}
}

func (p *ProducerHandlers) getAllProducers() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (p *ProducerHandlers) createProducer() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func (p *ProducerHandlers) getProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProducerHandlers) updateProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProducerHandlers) deleteProducerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterProducerRoutes(u usecase.Producer, g *echo.Group) {
	a := NewProducerHandlers(u)
	g.GET("", a.getAllProducers())
	g.POST("", a.createProducer())
	g.GET("/:id", a.getProducerById())
	g.PUT("/:id", a.updateProducerById())
	g.DELETE("/:id", a.deleteProducerById())
}
