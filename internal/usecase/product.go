package usecase

import (
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type product struct {
}

func NewProduct() Product {
	return &product{}
}

func (p *product) GetAll(filter *entity.ProductFilter) ([]*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) GetById(id uuid.UUID) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) Create(product entity.Product) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) UpdateById(id uuid.UUID, product entity.Product) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) DeleteById(id uuid.UUID) error {
	return ErrNotImplemented
}
