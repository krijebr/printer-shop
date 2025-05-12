package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type UserRepoPg struct {
	db *sql.DB
}

func NewUserRepoPg(db *sql.DB) User {
	return &UserRepoPg{
		db: db,
	}
}
func (u *UserRepoPg) GetAll(ctx context.Context, filter *entity.UserFilter) (allUsers []*entity.User, err error) {
	return nil, nil
}
func (u *UserRepoPg) GetById(ctx context.Context, id uuid.UUID) (user *entity.User, err error) {
	return nil, nil
}
func (u *UserRepoPg) Create(ctx context.Context, user entity.User) (createdUser *entity.User, err error) {
	return nil, nil
}
func (u *UserRepoPg) UpdateById(ctx context.Context, id uuid.UUID, user entity.User) (updatedUser *entity.User, err error) {
	return nil, nil
}
func (u *UserRepoPg) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	return nil
}
func (u *UserRepoPg) GetByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	return nil, nil
}
