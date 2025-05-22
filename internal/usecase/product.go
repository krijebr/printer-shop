package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type product struct {
	repo            repo.Product
	producerUseCase Producer
}

func NewProduct(r repo.Product, p Producer) Product {
	return &product{
		repo:            r,
		producerUseCase: p,
	}
}

func (p *product) GetAll(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, error) {
	return p.repo.GetAll(ctx, filter)
}
func (p *product) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) Create(ctx context.Context, product entity.Product) (*entity.Product, error) {
	_, err := p.producerUseCase.GetById(ctx, product.Producer.Id)
	if err != nil {
		return nil, err
	}
	product.Id = uuid.New()
	product.CreatedAt = time.Now()
	err = p.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	newProduct, err := p.repo.GetById(ctx, product.Id)
	if err != nil {
		return nil, err
	}
	return newProduct, nil
}
func (p *product) UpdateById(ctx context.Context, id uuid.UUID, product entity.Product) (*entity.Product, error) {
	return nil, ErrNotImplemented
}
func (p *product) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
