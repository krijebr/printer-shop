package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", gethandler)
	e.POST("/", posthandler)
	e.PUT("/", puthandler)
	e.DELETE("/", deletehandler)

	e.Logger.Fatal(e.Start(":8000"))
}

func gethandler(c echo.Context) error {
	c.String(http.StatusOK, "Ответ на GET запрос")
	return nil
}

func posthandler(c echo.Context) error {
	c.String(http.StatusOK, "Ответ на POST запрос")
	return nil
}

func puthandler(c echo.Context) error {
	c.String(http.StatusOK, "Ответ на PUT запрос")
	return nil
}

func deletehandler(c echo.Context) error {
	c.String(http.StatusOK, "Ответ на DELETE запрос")
	return nil
}
