package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	var ctx = context.Background()

	err := rdb.Set(ctx, "message", "example", time.Second*50).Err()
	if err != nil {
		log.Println(err)
	}
	result, _ := rdb.Get(ctx, "message").Result()

	log.Println(result)

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
