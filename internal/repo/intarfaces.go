package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type User interface {
	GetAll(ctx context.Context, filter *entity.UserFilter) (allUsers []*entity.User, err error)
	GetById(ctx context.Context, id uuid.UUID) (user *entity.User, err error)
	Create(ctx context.Context, user entity.User) (err error)
	Update(ctx context.Context, user entity.User) (err error)
	DeleteById(ctx context.Context, id uuid.UUID) (err error)
	GetByEmail(ctx context.Context, email string) (user *entity.User, err error)
}
