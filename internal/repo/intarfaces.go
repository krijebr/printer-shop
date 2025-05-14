package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type Token interface {
	SetToken(ctx context.Context, userId uuid.UUID, secret string, ttl time.Duration) (err error)
	SetRefreshToken(ctx context.Context, userId uuid.UUID, secret string, ttl time.Duration) (err error)
	GetTokenByUserId(ctx context.Context, userId uuid.UUID) (token string, err error)
	GetRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) (refrshToken string, err error)
	DeleteToken(ctx context.Context, userId uuid.UUID) (err error)
	DeleteRefreshToken(ctx context.Context, userId uuid.UUID) (err error)
}
type User interface {
	GetAll(ctx context.Context, filter *entity.UserFilter) (allUsers []*entity.User, err error)
	GetById(ctx context.Context, id uuid.UUID) (user *entity.User, err error)
	Create(ctx context.Context, user entity.User) (err error)
	Update(ctx context.Context, user entity.User) (err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
	GetByEmail(ctx context.Context, email string) (user *entity.User, err error)
}

type Producer interface {
	GetAll(ctx context.Context) (allProducers []*entity.Producer, err error)
	GetById(ctx context.Context, id uuid.UUID) (producer *entity.Producer, err error)
	Create(ctx context.Context, producer entity.Producer) (createdProducer *entity.Producer, err error)
	Update(ctx context.Context, producer entity.Producer) (updatedProducer *entity.Producer, err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
}
type Product interface {
	GetAll(ctx context.Context, filter *entity.ProductFilter) (allProducts []*entity.Product, err error)
	GetById(ctx context.Context, id uuid.UUID) (product *entity.Product, err error)
	Create(ctx context.Context, product entity.Product) (createdProduct *entity.Product, err error)
	Update(ctx context.Context, product entity.Product) (updatedProduct *entity.Product, err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
}
type Cart interface {
	GetAllProducts(ctx context.Context) (allProducts []*entity.ProductInCart, err error)
	AddProduct(ctx context.Context, productId uuid.UUID, count int) (err error)
	UpdateCount(ctx context.Context, productId uuid.UUID, count int) (err error)
}
