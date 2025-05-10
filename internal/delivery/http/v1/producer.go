package v1

import (
	"net/http"

	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProducerHandlers struct {
	usecase usecase.Producer
}

func NewProducerHandlers(uc usecase.Producer) *ProducerHandlers {
	return &ProducerHandlers{usecase: uc}
}

func (h *ProducerHandlers) getAllProducers(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *ProducerHandlers) createProducer(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}

func (h *ProducerHandlers) getProducerById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *ProducerHandlers) updateProducerById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
func (h *ProducerHandlers) deleteProducerById(c echo.Context) error {
	c.String(http.StatusNotImplemented, "Not Implemented")
	return nil
}
