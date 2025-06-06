package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
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
	confPath         string        = "./config/config.json"
	migratePath      string        = "file://./migrations"
	_defaultAttempts int           = 5
	_defaultTimeout  time.Duration = 5 * time.Second
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
		slog.Error("migrate up error", slog.Any("error", err))
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
	cartRepo := repo.NewCartRepoPg(db)
	orderRepo := repo.NewOrderRepoPg(db)

	producerUseCase := usecase.NewProducer(producerRepo, productRepo)
	authUseCase := usecase.NewAuth(userRepo, tokenRepo, time.Duration(cfg.Security.TokenTTL), time.Duration(cfg.Security.RefreshTokenTTL), cfg.Security.HashSalt)
	userUseCase := usecase.NewUser(userRepo, authUseCase)
	u := usecase.NewUseCases(authUseCase, usecase.NewCart(cartRepo, productRepo), usecase.NewOrder(orderRepo), producerUseCase, usecase.NewProduct(productRepo, producerRepo, cartRepo), userUseCase)
	r := http.CreateNewEchoServer(u)
	slog.Info("starting http server", slog.Int("port", cfg.HttpServer.Port))
	err = r.Start(fmt.Sprintf(":%d", cfg.HttpServer.Port))
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
		slog.Info(fmt.Sprintf("trying to initialize database, attempts left: %d", attempts))
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			_, err = db.Exec("select 1")
			if err == nil {
				break
			}
		}
		time.Sleep(_defaultTimeout)
		attempts--
	}
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
	m, err := migrate.NewWithDatabaseInstance(migratePath, "postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if err == migrate.ErrNoChange {
		slog.Info("migrate no changes")
		return nil
	}
	slog.Info("migrate up success")
	return nil
}
func initRedis(cfg *config.Redis) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ctx := context.Background()
	var (
		attempts = _defaultAttempts
		client   *redis.Client
		err      error
	)
	for attempts > 0 {
		slog.Info(fmt.Sprintf("trying to initialize redis, attempts left: %d", attempts))
		client = redis.NewClient(&redis.Options{Addr: addr, Password: cfg.Password, DB: cfg.DB})
		_, err = client.Ping(ctx).Result()
		if err == nil {
			break
		}
		time.Sleep(_defaultTimeout)
		attempts--
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}
