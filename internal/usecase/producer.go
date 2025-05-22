package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type producer struct {
	repo        repo.Producer
	repoProduct repo.Product
}

func NewProducer(r repo.Producer, p repo.Product) Producer {
	return &producer{
		repo:        r,
		repoProduct: p,
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
func (p *producer) Update(ctx context.Context, producer entity.Producer) (*entity.Producer, error) {
	_, err := p.repo.GetById(ctx, producer.Id)
	if err != nil {
		switch {
		case err == repo.ErrProducerNotFound:
			return nil, ErrProducerNotFound
		default:
			return nil, err
		}
	}
	err = p.repo.Update(ctx, producer)
	if err != nil {
		return nil, err
	}
	updatedProducer, err := p.repo.GetById(ctx, producer.Id)
	if err != nil {
		return nil, err
	}
	return updatedProducer, nil
}
func (p *producer) DeleteById(ctx context.Context, id uuid.UUID) error {
	filter := &entity.ProductFilter{
		ProducerId: &id,
	}
	products, err := p.repoProduct.GetAll(ctx, filter)
	if err != nil {
		return err
	}
	if len(products) != 0 {
		return ErrProducerUsed
	}
	_, err = p.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case err == repo.ErrProducerNotFound:
			return ErrProducerNotFound
		default:
			return err
		}
	}
	err = p.repo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
