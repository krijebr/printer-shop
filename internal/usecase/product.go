package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type product struct {
}

func NewProduct() Product {
	return &product{}
}

func (p *product) GetAll(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) Create(ctx context.Context, product entity.Product) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) UpdateById(ctx context.Context, id uuid.UUID, product entity.Product) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
