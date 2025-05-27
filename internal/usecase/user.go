package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type user struct {
	repo        repo.User
	authUseCase Auth
}

func NewUser(r repo.User, authUseCase Auth) User {
	return &user{
		repo:        r,
		authUseCase: authUseCase,
	}
}

func (u *user) GetAll(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return u.repo.GetAll(ctx, filter)
}
func (u *user) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetById(ctx, id)
}

func (u *user) Update(ctx context.Context, user entity.User) (*entity.User, error) {
	_, err := u.repo.GetById(ctx, user.Id)
	if err != nil {
		switch {
		case err == repo.ErrUserNotFound:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	if user.PasswordHash != "" {
		user.PasswordHash = u.authUseCase.HashPassword(user.PasswordHash)
	}
	err = u.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	updatedUser, err := u.repo.GetById(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
func (u *user) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
