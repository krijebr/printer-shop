package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type cart struct {
	repo        repo.Cart
	repoProduct repo.Product
}

func NewCart(r repo.Cart, p repo.Product) Cart {
	return &cart{
		repo:        r,
		repoProduct: p,
	}
}

func (c *cart) GetAllProducts(ctx context.Context, userId uuid.UUID) ([]*entity.ProductInCart, error) {
	ProductsInCart, err := c.repo.GetAllProducts(ctx, userId)
	if err != nil {
		return nil, err
	}
	return ProductsInCart, nil
}

func (c *cart) AddProduct(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) error {
	product, err := c.repoProduct.GetById(ctx, productId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			return ErrProductNotFound
		default:
			return err
		}
	}
	if product.Status == entity.ProductStatusHidden && count > 0 {
		return ErrProductIsHidden
	}
	existingCount, err := c.repo.GetProductCountById(ctx, userId, productId)
	if err != nil {
		return err
	}

	if existingCount == 0 {
		if count == 0 {
			return nil
		}
		err = c.repo.AddProduct(ctx, userId, productId, count)
		if err != nil {
			return err
		}
	}
	if count == 0 {
		err = c.repo.DeleteByProductId(ctx, userId, productId)
		if err != nil {
			return err
		}
	}
	err = c.repo.UpdateCount(ctx, userId, productId, count)
	if err != nil {
		return err
	}

	return nil
}

func (c *cart) UpdateCount(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) error {
	return nil
}
