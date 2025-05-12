package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type order struct {
}

func NewOrder() Order {
	return &order{}
}
func (o *order) Create(ctx context.Context, userId uuid.UUID) (err error) {
	return ErrNotImplemented
}
func (o *order) GetAll(ctx context.Context, filter *entity.OrderFilter) ([]*entity.Order, error) {
	return nil, ErrNotImplemented
}
func (o *order) GetById(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	return nil, ErrNotImplemented
}
func (o *order) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
func (o *order) UpdateById(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	return nil, ErrNotImplemented
}
