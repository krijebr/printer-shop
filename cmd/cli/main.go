package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "embed"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/krijebr/printer-shop/internal/config"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/urfave/cli/v2"
)

type ActionsCli struct {
	authUseCase     usecase.Auth
	userUseCase     usecase.User
	producerUseCase usecase.Producer
	productUseCase  usecase.Product
}

func NewActionsCli(a usecase.Auth, u usecase.User, p usecase.Producer, pr usecase.Product) *ActionsCli {
	return &ActionsCli{
		authUseCase:     a,
		userUseCase:     u,
		producerUseCase: p,
		productUseCase:  pr,
	}
}

const (
	defaultPath      string        = "E:/GO_projects/printer-shop/config/config.json"
	demoDataPath     string        = "E:/GO_projects/printer-shop/demo/demo-data.json"
	_defaultAttempts int           = 5
	_defaultTimeout  time.Duration = 5 * time.Second
)

//go:embed demo-data.json
var data []byte

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.InitConfigFromJson(getConfigPath(defaultPath))
	if err != nil {
		fmt.Println("initialization error", err)
		return
	}

	db, err := initDB(ctx, &cfg.Postgres)
	if err != nil {
		fmt.Println("database intialization error", err)
		return
	}

	userRepo := repo.NewUserRepoPg(db)
	producerRepo := repo.NewProducerRepoPg(db)
	productRepo := repo.NewProductRepoPg(db)
	tokenRepo := repo.NewTokenRedis(nil)
	cartRepo := repo.NewCartRepoPg(db)
	orderRepo := repo.NewOrderRepoPg(db)

	producerUseCase := usecase.NewProducer(producerRepo, productRepo)

	authUseCase := usecase.NewAuth(
		userRepo,
		tokenRepo,
		time.Duration(cfg.Security.TokenTTL),
		time.Duration(cfg.Security.RefreshTokenTTL),
		cfg.Security.HashSalt)
	userUseCase := usecase.NewUser(userRepo, cartRepo, orderRepo, authUseCase)
	productUseCase := usecase.NewProduct(productRepo, producerRepo, cartRepo, orderRepo)
	actionsCli := NewActionsCli(authUseCase, userUseCase, producerUseCase, productUseCase)

	app := &cli.App{
		Name:  "cli",
		Usage: "cli application of printer shop",
		Commands: []*cli.Command{
			{
				Name:   "create-admin",
				Usage:  "creates a new admin",
				Action: actionsCli.CreateAdmin(),
			},
			{
				Name:   "add-demo-data",
				Usage:  "fills the database with demo data",
				Action: actionsCli.AddDemoData(),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}

func (a *ActionsCli) CreateAdmin() cli.ActionFunc {
	return func(c *cli.Context) error {

		firstName := c.Args().Get(0)
		lastName := c.Args().Get(1)
		email := c.Args().Get(2)
		password := c.Args().Get(3)
		validate := validator.New()

		err := validate.Var(firstName, "required,max=25,min=3")
		if err != nil {
			fmt.Println("validation error")
			return err
		}
		err = validate.Var(lastName, "required,max=25,min=3")
		if err != nil {
			fmt.Println("validation error")
			return err
		}
		err = validate.Var(email, "required,email")
		if err != nil {
			fmt.Println("validation error")
			return err
		}
		err = validate.Var(password, "required,max=60,min=8")
		if err != nil {
			fmt.Println("validation error")
			return err
		}

		user := entity.User{
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			PasswordHash: password,
		}
		newUser, err := a.authUseCase.Register(c.Context, user)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrEmailAlreadyExists):
				fmt.Println("user with this email already exists")
				return nil
			default:
				fmt.Println("user creation error")
				return err
			}
		}
		newUser.Role = entity.UserRoleAdmin
		newUser.PasswordHash = ""
		newUser, err = a.userUseCase.Update(c.Context, *newUser)
		if err != nil {
			fmt.Println("user updating error")
			return err
		}
		fmt.Printf("new admin created. Email: %s, FirstName: %s, LastName: %s\n", newUser.Email, newUser.FirstName, newUser.LastName)
		return nil
	}
}

func (a *ActionsCli) AddDemoData() cli.ActionFunc {
	type Product struct {
		Producer entity.Producer  `json:"producer"`
		Products []entity.Product `json:"products"`
	}
	return func(c *cli.Context) error {
		demoData := []Product{}
		err := json.Unmarshal(data, &demoData)
		if err != nil {
			fmt.Println("datafile parsing error")
			return err
		}
		for i := range demoData {
			newProducer, err := a.producerUseCase.Create(c.Context, demoData[i].Producer)
			if err != nil {
				fmt.Println("producer creation error")
				return err
			}
			for j := range demoData[i].Products {
				demoData[i].Products[j].Producer = newProducer
				_, err = a.productUseCase.Create(c.Context, demoData[i].Products[j])
				if err != nil {
					fmt.Println("product creation error")
					return err
				}
			}
		}
		fmt.Println("demo data successfully added to database")
		return nil
	}
}

func initDB(ctx context.Context, cfg *config.Postgres) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBName)
	var (
		attempts = _defaultAttempts
		err      error
		db       *sql.DB
	)
	for attempts > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		fmt.Printf("trying to initialize database, attempts left: %d\n", attempts)
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			_, err = db.ExecContext(ctx, "select 1")
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

func getConfigPath(defaultPath string) string {
	if os.Getenv("CONFIG_PATH") != "" {
		return os.Getenv("CONFIG_PATH")
	}
	return defaultPath
}
