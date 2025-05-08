package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (token string, refreshToken string, err error)
	ValidateToken(ctx context.Context, token string) (user *entity.User, err error)
	RefreshToken(ctx context.Context, refreshToken string) (token string, newRefreshToken string, err error)
}

type User interface {
	GetAll() (allUsers []*entity.User, err error)
	GetById(id uuid.UUID) (user *entity.User, err error)
	Create(user entity.User) (createdUser *entity.User, err error)
	UpdateById(id uuid.UUID, user entity.User) (updatedUser *entity.User, err error)
	DeleteById(id uuid.UUID) (err error)
}

type Product interface {
	GetAll() (allProducts []*entity.ProductWithProducer, err error)
	GetById(id uuid.UUID) (product *entity.ProductWithProducer, err error)
	Create(product entity.Product) (createdProduct *entity.Product, err error)
	UpdateById(id uuid.UUID, product entity.Product) (updatedProduct *entity.ProductWithProducer, err error)
	DeleteById(id uuid.UUID) (err error)
}

type Producer interface {
	GetAll() (allProducers []*entity.Producer, err error)
	GetById(id uuid.UUID) (producer *entity.Producer, err error)
	Create(producer entity.Producer) (createdProducer *entity.Producer, err error)
	UpdateById(id uuid.UUID, producer entity.Producer) (updatedProducer *entity.Producer, err error)
	DeleteById(id uuid.UUID) (err error)
}

type Cart interface {
	GetAllProducts() (allProducts []*entity.ProductInCart, err error)
	AddProduct(productId uuid.UUID, count int) (err error)
	UpdateCount(productId uuid.UUID) (err error)
}

type Order interface {
	Place([]*entity.ProductForOrder) (err error)
	GetAll() (allOrders []*entity.OrderWithProducts, err error)
	GetById(id uuid.UUID) (allProducts entity.OrderWithProducts, err error)
	UpdateById(id uuid.UUID) (order entity.OrderWithProducts, err error)
	ChangeStatus(productId uuid.UUID, status entity.ProductStatus) (order entity.OrderWithProducts, err error)
}
