package usecase

import (
	"context"

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
	return nil
}
func (c *cart) UpdateCount(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) error {
	return nil
}
