package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type Auth interface {
	Register(ctx context.Context, user entity.User) (createdUser *entity.User, err error)
	Login(ctx context.Context, email, password string) (token string, refreshToken string, err error)
	ValidateToken(ctx context.Context, token string) (user *entity.User, err error)
	RefreshToken(ctx context.Context, refreshToken string) (token string, newRefreshToken string, err error)
	ValidatePassword(password, hash string) (result bool)
	HashPassword(password string) (hashPassword string)
}

type User interface {
	GetAll(ctx context.Context, filter *entity.UserFilter) (allUsers []*entity.User, err error)
	GetById(ctx context.Context, id uuid.UUID) (user *entity.User, err error)
	Update(ctx context.Context, user entity.User) (updatedUser *entity.User, err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
}

type Product interface {
	GetAll(ctx context.Context, filter *entity.ProductFilter) (allProducts []*entity.Product, err error)
	GetById(ctx context.Context, id uuid.UUID) (product *entity.Product, err error)
	Create(ctx context.Context, product entity.Product) (createdProduct *entity.Product, err error)
	Update(ctx context.Context, product entity.Product) (updatedProduct *entity.Product, err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
}

type Producer interface {
	GetAll(ctx context.Context) (allProducers []*entity.Producer, err error)
	GetById(ctx context.Context, id uuid.UUID) (producer *entity.Producer, err error)
	Create(ctx context.Context, producer entity.Producer) (createdProducer *entity.Producer, err error)
	Update(ctx context.Context, producer entity.Producer) (updatedProducer *entity.Producer, err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
}

type Cart interface {
	GetAllProducts(ctx context.Context, userId uuid.UUID) (allProducts []*entity.ProductInCart, err error)
	AddProduct(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) (err error)
	UpdateCount(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) (err error)
}

type Order interface {
	Create(ctx context.Context, userId uuid.UUID) (order *entity.Order, err error)
	GetAll(ctx context.Context, filter *entity.OrderFilter) (allOrders []*entity.Order, err error)
	GetById(ctx context.Context, id uuid.UUID) (order *entity.Order, err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
	UpdateById(ctx context.Context, id uuid.UUID) (order *entity.Order, err error)
}
