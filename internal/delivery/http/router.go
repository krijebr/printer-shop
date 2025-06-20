package http

import (
	"github.com/krijebr/printer-shop/internal/config"
	"github.com/krijebr/printer-shop/internal/delivery/http/middlewares"
	v1 "github.com/krijebr/printer-shop/internal/delivery/http/v1"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

func CreateNewEchoServer(u *usecase.UseCases, r *config.RoleConf, baseUrl string) *echo.Echo {
	authMw := middlewares.NewAuthMiddleware(u, r, baseUrl)
	server := echo.New()
	server.HideBanner = true
	server.GET("health", HealthCheck())
	g := server.Group(baseUrl)
	v1.RegisterAuthRoutes(u.Auth, g.Group("auth"))
	v1.RegisterCartRoutes(u.Cart, g.Group("cart", authMw.Handle))
	v1.RegisterOrderRoutes(u.Order, g.Group("orders", authMw.Handle))
	v1.RegisterProducerRoutes(u.Producer, g.Group("producers", authMw.Handle))
	v1.RegisterProductRoutes(u.Product, g.Group("products", authMw.Handle))
	v1.RegisterProfileRoutes(u.User, g.Group("profile", authMw.Handle))
	v1.RegisterUserRoutes(u.User, g.Group("users", authMw.Handle))
	return server
}
