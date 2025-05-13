package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type user struct {
	repo repo.User
}

func NewUser(r repo.User) User {
	return &user{repo: r}
}

func (u *user) GetAll(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	users, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	return users, err
}
func (u *user) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) Register(ctx context.Context, user entity.User) (*entity.User, error) {
	user.Id = uuid.New()
	user.CreatedAt = time.Now()
	user.Status = entity.UserStatusActive
	user.Role = entity.UserRoleCustomer
	err := u.repo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	newUser, err := u.repo.GetById(ctx, user.Id)
	return newUser, nil
}
func (u *user) UpdateById(ctx context.Context, id uuid.UUID, user entity.User) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
