package main

import (
	"log"

	"github.com/krijebr/printer-shop/internal/delivery/http"
	"github.com/krijebr/printer-shop/internal/usecase"
)

func main() {

	/*rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	var ctx = context.Background()

	err := rdb.Set(ctx, "message", "example", time.Second*50).Err()
	if err != nil {
		log.Println(err)
	}
	result, _ := rdb.Get(ctx, "message").Result()

	log.Println(result)*/
	u := usecase.NewUseCases(usecase.NewAuth(), usecase.NewCart(), usecase.NewOrder(), usecase.NewProducer(), usecase.NewProduct(), usecase.NewUser())
	r := http.CreateNewEchoServer(u)

	err := r.Start(":8000")
	if err != nil {
		log.Println("Ошибка запуска сервера")
	}
}
