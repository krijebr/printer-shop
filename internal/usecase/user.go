package usecase

import (
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type user struct {
}

func NewUser() User {
	return &user{}
}

func (u *user) GetAll() ([]*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) GetById(id uuid.UUID) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) Create(user entity.User) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) UpdateById(id uuid.UUID, user entity.User) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (u *user) DeleteById(id uuid.UUID) error {

	return ErrNotImplemented
}
