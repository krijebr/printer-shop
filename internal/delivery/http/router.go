package http

import (
	"github.com/krijebr/printer-shop/internal/delivery/http/middlewares"
	v1 "github.com/krijebr/printer-shop/internal/delivery/http/v1"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

func CreateNewEchoServer(u *usecase.UseCases) *echo.Echo {
	authMw := middlewares.NewAuthMiddleware(u)
	server := echo.New()
	server.HideBanner = true
	g := server.Group("/api/v1/")
	v1.RegisterAuthRoutes(u.Auth, g.Group("auth"))
	v1.RegisterCartRoutes(u.Cart, g.Group("cart", authMw.Handle))
	v1.RegisterOrderRoutes(u.Order, g.Group("orders", authMw.Handle))
	v1.RegisterProducerRoutes(u.Producer, g.Group("producers"))
	v1.RegisterProductRoutes(u.Product, g.Group("products"))
	v1.RegisterProfileRoutes(u.User, g.Group("profile", authMw.Handle))
	v1.RegisterUserRoutes(u.User, g.Group("users"))
	return server
}
