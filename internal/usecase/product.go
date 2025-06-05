package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type product struct {
	repo         repo.Product
	repoProducer repo.Producer
	repoCart     repo.Cart
}

func NewProduct(r repo.Product, p repo.Producer, c repo.Cart) Product {
	return &product{
		repo:         r,
		repoProducer: p,
		repoCart:     c,
	}
}

func (p *product) GetAll(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, error) {
	return p.repo.GetAll(ctx, filter)
}
func (p *product) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := p.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case err == repo.ErrProductNotFound:
			return nil, ErrProductNotFound
		default:
			return nil, err
		}
	}
	return product, nil
}
func (p *product) Create(ctx context.Context, product entity.Product) (*entity.Product, error) {
	_, err := p.repoProducer.GetById(ctx, product.Producer.Id)
	if err != nil {
		switch {
		case err == repo.ErrProducerNotFound:
			return nil, ErrProducerNotFound
		default:
			return nil, err
		}
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
func (p *product) Update(ctx context.Context, product entity.Product) (*entity.Product, error) {
	_, err := p.repo.GetById(ctx, product.Id)
	if err != nil {
		switch {
		case err == repo.ErrProductNotFound:
			return nil, ErrProductNotFound
		default:
			return nil, err
		}
	}
	if product.Producer.Id != uuid.Nil {
		_, err = p.repoProducer.GetById(ctx, product.Producer.Id)
		if err != nil {
			switch {
			case err == repo.ErrProducerNotFound:
				return nil, ErrProducerNotFound
			default:
				return nil, err
			}
		}
	}
	err = p.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}
	updatedProduct, err := p.repo.GetById(ctx, product.Id)
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}
func (p *product) DeleteById(ctx context.Context, id uuid.UUID) error {
	productExists, err := p.repoCart.CheckIfExistsById(ctx, id)
	if err != nil {
		return err
	}
	if productExists {
		return ErrProductIsUsed
	}
	_, err = p.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case err == repo.ErrProductNotFound:
			return ErrProductNotFound
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
