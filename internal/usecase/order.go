package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type order struct {
	repo     repo.Order
	repoCart repo.Cart
}

func NewOrder(r repo.Order, c repo.Cart) Order {
	return &order{
		repo:     r,
		repoCart: c,
	}
}
func (o *order) Create(ctx context.Context, userId uuid.UUID) (*entity.Order, error) {
	productsInCart, err := o.repoCart.GetAllProducts(ctx, userId)
	if len(productsInCart) == 0 {
		return nil, ErrCartIsEmpty
	}

	order := &entity.Order{}
	order.Id = uuid.New()
	order.UserId = userId
	order.Status = entity.OrderStatusNew
	order.CreatedAt = time.Now()
	order.Products = productsInCart
	err = o.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	newOrder, err := o.repo.GetById(ctx, order.Id)
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}
func (o *order) GetAll(ctx context.Context, filter *entity.OrderFilter) ([]*entity.Order, error) {
	order, err := o.repo.GetAll(ctx, filter)
	if err != nil {
		switch {
		case err == repo.ErrOrderNotFound:
			return nil, ErrOrderNotFound
		default:
			return nil, err
		}
	}
	return order, nil
}
func (o *order) GetById(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	return o.repo.GetById(ctx, id)
}
func (o *order) DeleteById(ctx context.Context, id uuid.UUID) error {
	return ErrNotImplemented
}
func (o *order) UpdateById(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	return nil, ErrNotImplemented
}
