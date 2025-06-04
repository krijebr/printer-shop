package v1

import (
	"log/slog"
	"net/http"

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
			return c.NoContent(http.StatusInternalServerError)
		}
		productsInCart, err := h.usecase.GetAllProducts(c.Request().Context(), userId)
		if err != nil {
			slog.Error("profile receiving error", slog.Any("error", err))
			return c.NoContent(http.StatusInternalServerError)
		}
		slog.Info("cart received")
		return c.JSON(http.StatusOK, productsInCart)
	}
}
func (h *CartHandlers) addProductToCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterCartRoutes(u usecase.Cart, g *echo.Group) {
	a := NewCartHandlers(u)
	g.GET("", a.getAllProductsInCart())
	g.POST("", a.addProductToCart())
}
