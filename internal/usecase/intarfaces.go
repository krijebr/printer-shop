package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (token string, refresh_token string, err error)
	ValidateToken(ctx context.Context, token string) (user *entity.User, err error)
	Refresh(ctx context.Context, refresh_token string) (token string, new_refresh_token string, err error)
}

type User interface {
	GetAll() (all_users []*entity.User, err error)
	GetById(id uuid.UUID) (user *entity.User, err error)
	Create(user entity.User) (created_user *entity.User, err error)
	UpdateById(id uuid.UUID, user entity.User) (updated_user *entity.User, err error)
	DeleteById(id uuid.UUID) (err error)
}

type Product interface {
	GetAll() (all_products []*entity.Product, err error)
	GetById(id uuid.UUID) (product *entity.Product, err error)
	Create(product entity.Product) (created_product *entity.Product, err error)
	UpdateById(id uuid.UUID, product entity.Product) (updated_product *entity.Product, err error)
	DeleteById(id uuid.UUID) (err error)
}

type Producer interface {
	GetAll() (all_producers []*entity.Producer, err error)
	GetById(id uuid.UUID) (producer *entity.Producer, err error)
	Create(producer entity.Producer) (created_producer *entity.Producer, err error)
	UpdateById(id uuid.UUID, producer entity.Producer) (updated_producer *entity.Producer, err error)
	DeleteById(id uuid.UUID) (err error)
}
