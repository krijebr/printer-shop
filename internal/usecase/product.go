package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type product struct {
	repo         repo.Product
	repoProducer repo.Producer
	repoCart     repo.Cart
	repoOrder    repo.Order
}

func NewProduct(r repo.Product, p repo.Producer, c repo.Cart, o repo.Order) Product {
	return &product{
		repo:         r,
		repoProducer: p,
		repoCart:     c,
		repoOrder:    o,
	}
}

func (p *product) GetAll(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, error) {
	return p.repo.GetAll(ctx, filter)
}
func (p *product) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	receivedProduct, err := p.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			return nil, ErrProductNotFound
		default:
			return nil, err
		}
	}
	return receivedProduct, nil
}
func (p *product) Create(ctx context.Context, productToCreate entity.Product) (*entity.Product, error) {
	_, err := p.repoProducer.GetById(ctx, productToCreate.Producer.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProducerNotFound):
			return nil, ErrProducerNotFound
		default:
			return nil, err
		}
	}
	productToCreate.Id = uuid.New()
	productToCreate.CreatedAt = time.Now()
	err = p.repo.Create(ctx, productToCreate)
	if err != nil {
		return nil, err
	}
	newProduct, err := p.repo.GetById(ctx, productToCreate.Id)
	if err != nil {
		return nil, err
	}
	return newProduct, nil
}
func (p *product) Update(ctx context.Context, productToUpdate entity.Product) (*entity.Product, error) {
	_, err := p.repo.GetById(ctx, productToUpdate.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			return nil, ErrProductNotFound
		default:
			return nil, err
		}
	}
	if productToUpdate.Producer.Id != uuid.Nil {
		_, err = p.repoProducer.GetById(ctx, productToUpdate.Producer.Id)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrProducerNotFound):
				return nil, ErrProducerNotFound
			default:
				return nil, err
			}
		}
	}
	err = p.repo.Update(ctx, productToUpdate)
	if err != nil {
		return nil, err
	}
	updatedProduct, err := p.repo.GetById(ctx, productToUpdate.Id)
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}
func (p *product) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := p.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			return ErrProductNotFound
		default:
			return err
		}
	}
	productExists, err := p.repoCart.CheckIfExistsById(ctx, id)
	if err != nil {
		return err
	}
	if productExists {
		return ErrProductIsUsed
	}
	productExists, err = p.repoOrder.CheckIfExistsByProductId(ctx, id)
	if err != nil {
		return err
	}
	if productExists {
		return ErrProductIsUsed
	}
	err = p.repo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
