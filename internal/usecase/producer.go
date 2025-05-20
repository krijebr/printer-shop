package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type producer struct {
	repo repo.Producer
}

func NewProducer(r repo.Producer) Producer {
	return &producer{
		repo: r,
	}
}
func (p *producer) GetAll(ctx context.Context) ([]*entity.Producer, error) {
	return p.repo.GetAll(ctx)
}
func (p *producer) GetById(ctx context.Context, id uuid.UUID) (*entity.Producer, error) {
	producer, err := p.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case err == repo.ErrProducerNotFound:
			return nil, ErrProducerNotFound
		default:
			return nil, err
		}
	}
	return producer, nil
}
func (p *producer) Create(ctx context.Context, producer entity.Producer) (*entity.Producer, error) {
	producer.Id = uuid.New()
	producer.CreatedAt = time.Now()
	err := p.repo.Create(ctx, producer)
	if err != nil {
		return nil, err
	}
	newProducer, err := p.repo.GetById(ctx, producer.Id)
	if err != nil {
		return nil, err
	}
	return newProducer, nil
}
func (p *producer) UpdateById(ictx context.Context, d uuid.UUID, producer entity.Producer) (*entity.Producer, error) {
	return nil, ErrNotImplemented
}
func (p *producer) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
