package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type user struct {
	repo        repo.User
	repoCart    repo.Cart
	repoOrder   repo.Order
	authUseCase Auth
}

func NewUser(r repo.User, c repo.Cart, o repo.Order, authUseCase Auth) User {
	return &user{
		repo:        r,
		repoCart:    c,
		repoOrder:   o,
		authUseCase: authUseCase,
	}
}

func (u *user) GetAll(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *user) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetById(ctx, id)
}

func (u *user) Update(ctx context.Context, userToUpdate entity.User) (*entity.User, error) {
	_, err := u.repo.GetById(ctx, userToUpdate.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	if userToUpdate.PasswordHash != "" {
		userToUpdate.PasswordHash = u.authUseCase.HashPassword(userToUpdate.PasswordHash)
	}
	err = u.repo.Update(ctx, userToUpdate)
	if err != nil {
		return nil, err
	}
	updatedUser, err := u.repo.GetById(ctx, userToUpdate.Id)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (u *user) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			return ErrUserNotFound
		default:
			return err
		}
	}
	filter := &entity.OrderFilter{
		UserId: &id,
	}
	orders, err := u.repoOrder.GetAll(ctx, filter)
	if err != nil {
		return err
	}
	productsInCart, err := u.repoCart.GetAllProducts(ctx, id)
	if err != nil {
		return err
	}
	if len(productsInCart) > 0 || len(orders) > 0 {
		return ErrUserIsUsed
	}
	err = u.repo.DeleteById(ctx, id)
	return nil
}
