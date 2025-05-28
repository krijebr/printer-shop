package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/krijebr/printer-shop/internal/config"
	"github.com/krijebr/printer-shop/internal/delivery/http"
	"github.com/krijebr/printer-shop/internal/repo"
	"github.com/krijebr/printer-shop/internal/usecase"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const confpath string = "../config/config.json"

func main() {

	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(th)
	slog.SetDefault(logger)
	slog.Info("starting app", slog.String("app-name", "printer shop"))

	cfg, err := config.InitConfigFromJson(confpath)
	if err != nil {
		slog.Error("initialization error", slog.Any("error", err))
		return
	}

	db, err := initDB(&cfg.Postgres)
	if err != nil {
		slog.Error("database intialization error", slog.Any("error", err))
		return
	}
	rdb, err := initRedis(&cfg.Redis)
	if err != nil {
		slog.Error("redis initialization error", slog.Any("error", err))
		return
	}
	userRepo := repo.NewUserRepoPg(db)
	producerRepo := repo.NewProducerRepoPg(db)
	productRepo := repo.NewProductRepoPg(db)
	tokenRepo := repo.NewTokenRedis(rdb)

	producerUseCase := usecase.NewProducer(producerRepo, productRepo)
	authUseCase := usecase.NewAuth(userRepo, tokenRepo, time.Duration(cfg.Security.TokenTTL), time.Duration(cfg.Security.RefreshTokenTTL), cfg.Security.HashSalt)
	userUseCase := usecase.NewUser(userRepo, authUseCase)
	u := usecase.NewUseCases(authUseCase, usecase.NewCart(), usecase.NewOrder(), producerUseCase, usecase.NewProduct(productRepo, producerRepo), userUseCase)
	r := http.CreateNewEchoServer(u)
	r.HideBanner = true
	port := strconv.Itoa(cfg.HttpServer.Port)
	slog.Info("starting http server", slog.Any("port", port))
	err = r.Start(":" + port)
	if err != nil {
		slog.Error("starting server error", slog.Any("error", err))
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
