package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type order struct {
	repo        repo.Order
	repoCart    repo.Cart
	repoProduct repo.Product
}

func NewOrder(r repo.Order, c repo.Cart, p repo.Product) Order {
	return &order{
		repo:        r,
		repoCart:    c,
		repoProduct: p,
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
	publishedProducts := []*entity.ProductInCart{}
	for _, product := range productsInCart {
		if product.Product.Status != entity.ProductStatusHidden {
			publishedProducts = append(publishedProducts, product)
		}
	}
	order.Products = publishedProducts

	err = o.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	err = o.repoCart.ClearCart(ctx, userId)
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
		case errors.Is(err, repo.ErrOrderNotFound):
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
	_, err := o.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			return ErrOrderNotFound
		default:
			return err
		}
	}
	err = o.repo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (o *order) UpdateById(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	_, err := o.repo.GetById(ctx, order.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			return nil, ErrOrderNotFound
		default:
			return nil, err
		}
	}
	if order.Products != nil {
		publishedProducts := []*entity.ProductInCart{}
		for i, newProduct := range order.Products {
			product, err := o.repoProduct.GetById(ctx, newProduct.Product.Id)
			if err != nil {
				switch {
				case errors.Is(err, repo.ErrProductNotFound):
					return nil, ErrProductNotFound
				default:
					return nil, err
				}
			}
			if product.Status != entity.ProductStatusHidden {
				publishedProducts = append(publishedProducts, newProduct)
			}
			order.Products[i].Product.Price = product.Price
		}
		order.Products = publishedProducts
	}
	err = o.repo.UpdateById(ctx, order)
	if err != nil {
		return nil, err
	}
	updatedOrder, err := o.repo.GetById(ctx, order.Id)
	if err != nil {
		return nil, err
	}
	return updatedOrder, ErrNotImplemented
}
