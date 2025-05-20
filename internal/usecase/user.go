package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type user struct {
	repo     repo.User
	hashSalt string
}

func NewUser(r repo.User, salt string) User {
	return &user{
		repo:     r,
		hashSalt: salt,
	}
}

func (u *user) hashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass + u.hashSalt))
	return fmt.Sprintf("%x", h.Sum(nil))
}
func (u *user) ValidatePassword(password, hash string) bool {
	return u.hashPassword(password) == hash
}

func (u *user) GetAll(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return u.repo.GetAll(ctx, filter)
}
func (u *user) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetById(ctx, id)
}
func (u *user) Register(ctx context.Context, user entity.User) (*entity.User, error) {
	someUser, err := u.repo.GetByEmail(ctx, user.Email)
	if err != nil && err != repo.ErrUserNotFound {
		return nil, err
	}
	if someUser != nil {
		return nil, ErrEmailAlreadyExists
	}
	user.Id = uuid.New()
	user.PasswordHash = u.hashPassword(user.PasswordHash)
	user.CreatedAt = time.Now()
	user.Status = entity.UserStatusActive
	user.Role = entity.UserRoleCustomer
	err = u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	newUser, err := u.repo.GetById(ctx, user.Id)

	return newUser, nil
}
func (u *user) Update(ctx context.Context, user entity.User) (*entity.User, error) {
	if user.PasswordHash != "" {
		user.PasswordHash = u.hashPassword(user.PasswordHash)
	}
	err := u.repo.Update(ctx, user)
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
