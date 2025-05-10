package http

import (
	"github.com/labstack/echo/v4"
)

func CreateRouter() *echo.Echo {
	myRouter := echo.New()
	return myRouter
}
