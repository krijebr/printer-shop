package http

import (
	v1 "github.com/krijebr/printer-shop/internal/delivery/http/v1"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

func CreateNewEchoServer(u *usecase.UseCases) *echo.Echo {
	server := echo.New()
	g := server.Group("/api/v1/")
	v1.RegisterAuthRoutes(u.Auth, g.Group("auth"))
	v1.RegisterCartRoutes(u.Cart, g.Group("cart"))
	v1.RegisterOrderRoutes(u.Order, g.Group("orders"))
	v1.RegisterProducerRoutes(u.Producer, g.Group("producers"))
	v1.RegisterProductRoutes(u.Product, g.Group("products"))
	v1.RegisterProfileRoutes(u.User, g.Group("profile"))
	v1.RegisterUserRoutes(u.User, g.Group("users"))
	return server
}
