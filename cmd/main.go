package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/krijebr/printer-shop/internal/config"
	"github.com/krijebr/printer-shop/internal/delivery/http"
	"github.com/krijebr/printer-shop/internal/repo"
	"github.com/krijebr/printer-shop/internal/usecase"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	confPath         string = "../config/config.json"
	_defaultAttempts        = 8
	_defaultTimeout         = 1000 * time.Millisecond
)

func main() {

	cfg, err := config.InitConfigFromJson(confPath)
	if err != nil {
		slog.Error("initialization error", slog.Any("error", err))
		return
	}

	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Logging.Level,
	})
	logger := slog.New(th)
	slog.SetDefault(logger)

	slog.Info("starting app", slog.String("app-name", "printer shop"))

	db, err := initDB(&cfg.Postgres)
	if err != nil {
		slog.Error("database intialization error", slog.Any("error", err))
		return
	}
	err = migration(db)

	if err != nil {
		slog.Error("Migrate: up error", slog.Any("error", err))
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
	port := strconv.Itoa(cfg.HttpServer.Port)
	slog.Info("starting http server", slog.Any("port", port))
	err = r.Start(":" + port)
	if err != nil {
		slog.Error("starting server error", slog.Any("error", err))
	}
}

func initDB(cfg *config.Postgres) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBName)

	var (
		attempts = _defaultAttempts
		err      error
		db       *sql.DB
	)
	for attempts > 0 {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			break
		}
		slog.Info(fmt.Sprintf("trying to initialize database, attempts left: %d", attempts))
		time.Sleep(_defaultTimeout)
		attempts--
	}
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("select 1")
	if err != nil {
		return nil, err
	}
	return db, nil
}
func migration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://E:/GO_projects/printer-shop/migrations", "postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	if err == migrate.ErrNoChange {
		slog.Info("Migrate: no changes")
		return nil
	}
	slog.Info("Migrate: up success")
	return nil
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
