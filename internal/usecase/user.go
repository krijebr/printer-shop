package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type user struct {
}

func NewUser() User {
	return &user{}
}

func (u *user) GetAll(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) UpdateById(ctx context.Context, id uuid.UUID, user entity.User) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) DeleteById(ctx context.Context, id uuid.UUID) error {

	return ErrNotImplemented
}
