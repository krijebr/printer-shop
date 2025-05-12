package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type producer struct {
}

func NewProducer() Producer {
	return &producer{}
}
func (p *producer) GetAll(ctx context.Context) ([]*entity.Producer, error) {
	return nil, ErrNotImplemented
}
func (p *producer) GetById(ctx context.Context, id uuid.UUID) (*entity.Producer, error) {
	return nil, ErrNotImplemented
}
func (p *producer) Create(ctx context.Context, producer entity.Producer) (*entity.Producer, error) {
	return nil, ErrNotImplemented
}
func (p *producer) UpdateById(ictx context.Context, d uuid.UUID, producer entity.Producer) (*entity.Producer, error) {
	return nil, ErrNotImplemented
}
func (p *producer) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
