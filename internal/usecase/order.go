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
	newOrder := &entity.Order{
		Id:        uuid.New(),
		UserId:    userId,
		Status:    entity.OrderStatusNew,
		CreatedAt: time.Now(),
	}
	publishedProducts := make([]*entity.ProductInCart, 0, len(productsInCart))
	for _, p := range productsInCart {
		if p.Product.Status != entity.ProductStatusHidden {
			publishedProducts = append(publishedProducts, p)
		}
	}
	newOrder.Products = publishedProducts

	err = o.repo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	err = o.repoCart.ClearCart(ctx, userId)
	if err != nil {
		return nil, err
	}
	createdOrder, err := o.repo.GetById(ctx, newOrder.Id)
	if err != nil {
		return nil, err
	}
	return createdOrder, nil
}

func (o *order) GetAll(ctx context.Context, filter *entity.OrderFilter) ([]*entity.Order, error) {
	orders, err := o.repo.GetAll(ctx, filter)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			return nil, ErrOrderNotFound
		default:
			return nil, err
		}
	}
	return orders, nil
}

func (o *order) GetById(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	orderToReceive, err := o.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			return nil, ErrOrderNotFound
		default:
			return nil, err
		}
	}
	return orderToReceive, nil
}

func (o *order) DeleteById(ctx context.Context, id uuid.UUID) error {
	orderToDelete, err := o.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			return ErrOrderNotFound
		default:
			return err
		}
	}
	if orderToDelete.Status != entity.OrderStatusNew {
		return ErrOrderCantBeDeleted
	}
	err = o.repo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (o *order) UpdateById(ctx context.Context, orderToUpdate *entity.Order) (*entity.Order, error) {
	existingOrder, err := o.repo.GetById(ctx, orderToUpdate.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrOrderNotFound):
			return nil, ErrOrderNotFound
		default:
			return nil, err
		}
	}
	if orderToUpdate.Products != nil {
		if existingOrder.Status != entity.OrderStatusNew {
			return nil, ErrOrderCantBeUpdated
		}
		publishedProducts := make([]*entity.ProductInCart, 0, len(orderToUpdate.Products))
		for _, newProduct := range orderToUpdate.Products {
			product, err := o.repoProduct.GetById(ctx, newProduct.Product.Id)
			if err != nil {
				switch {
				case errors.Is(err, repo.ErrProductNotFound):
					return nil, ErrProductNotFound
				default:
					return nil, err
				}
			}
			newProduct.Product.Price = product.Price
			if product.Status != entity.ProductStatusHidden {
				publishedProducts = append(publishedProducts, newProduct)
			}

		}
		orderToUpdate.Products = publishedProducts
	}
	err = o.repo.UpdateById(ctx, orderToUpdate)
	if err != nil {
		return nil, err
	}
	updatedOrder, err := o.repo.GetById(ctx, orderToUpdate.Id)
	if err != nil {
		return nil, err
	}
	return updatedOrder, nil
}
