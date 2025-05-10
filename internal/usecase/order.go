package usecase

import (
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type order struct {
}

func NewOrder() Order {
	return &order{}
}
func (o *order) Create(userId uuid.UUID) (err error) {
	return ErrNotImplemented
}
func (o *order) GetAll(filter *entity.OrderFilter) ([]*entity.Order, error) {
	return nil, ErrNotImplemented
}
func (o *order) GetById(id uuid.UUID) (*entity.Order, error) {
	return nil, ErrNotImplemented
}
func (o *order) DeleteById(id uuid.UUID) error {
	return ErrNotImplemented
}
func (o *order) UpdateById(id uuid.UUID) (*entity.Order, error) {
	return nil, ErrNotImplemented
}
