package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/krijebr/printer-shop/internal/config"
	"github.com/krijebr/printer-shop/internal/delivery/http"
	"github.com/krijebr/printer-shop/internal/repo"
	"github.com/krijebr/printer-shop/internal/usecase"
	_ "github.com/lib/pq"
)

const confpath string = "../config/config.json"

func main() {

	/*rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	var ctx = context.Background()

	err := rdb.Set(ctx, "message", "example", time.Second*50).Err()
	if err != nil {
		log.Println(err)
	}
	result, _ := rdb.Get(ctx, "message").Result()

	log.Println(result)*/

	log.Println("starting app")

	cfg, err := config.InitConfigFromJson(confpath)
	if err != nil {
		log.Println("Ошибка инициализации", err)
		return
	}
	db, err := initDB(&cfg.Postgres)
	if err != nil {
		log.Println("Ошибка инициализации базы данных", err)
		return
	}
	userRepo := repo.NewUserRepoPg(db)
	u := usecase.NewUseCases(usecase.NewAuth(), usecase.NewCart(), usecase.NewOrder(), usecase.NewProducer(), usecase.NewProduct(), usecase.NewUser(userRepo))
	r := http.CreateNewEchoServer(u)

	err = r.Start(":8000")
	if err != nil {
		log.Println("Ошибка запуска сервера")
	}
}

func initDB(cfg *config.Postgres) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("select 1")
	if err != nil {
		return nil, err
	}
	return db, nil
}
