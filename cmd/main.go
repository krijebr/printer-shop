package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/krijebr/printer-shop/internal/config"
	"github.com/krijebr/printer-shop/internal/delivery/http"
	"github.com/krijebr/printer-shop/internal/repo"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

const confpath string = "../config/config.json"

func main() {

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
	rdb, err := initRedis(&cfg.Redis)
	if err != nil {
		log.Println("Ошибка инициализации redis", err)
		return
	}
	userRepo := repo.NewUserRepoPg(db)
	tokenRepo := repo.NewTokenRedis(rdb)
	u := usecase.NewUseCases(usecase.NewAuth(userRepo, tokenRepo, time.Duration(cfg.Security.TokenTTL), time.Duration(cfg.Security.RefreshTokenTTL),
		cfg.Security.HashSalt), usecase.NewCart(), usecase.NewOrder(), usecase.NewProducer(), usecase.NewProduct(), usecase.NewUser(userRepo))
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
func initRedis(cfg *config.Redis) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: addr, Password: cfg.Password, DB: cfg.DB})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
