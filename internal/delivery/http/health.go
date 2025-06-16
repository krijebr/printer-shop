package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheck() echo.HandlerFunc {
	type healthCheckResponse struct {
		Health string `json:"health"`
	}
	return func(c echo.Context) error {
		healthTrue := healthCheckResponse{
			Health: "true",
		}
		return c.JSON(http.StatusOK, healthTrue)
	}
}
