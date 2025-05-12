package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type cart struct {
}

func NewCart() Cart {
	return &cart{}
}
func (c *cart) GetAllProducts(ctx context.Context) (allProducts []*entity.ProductInCart, err error) {
	return nil, err
}
func (c *cart) AddProduct(ctx context.Context, productId uuid.UUID, count int) (err error) {
	return err
}
func (c *cart) UpdateCount(ctx context.Context, productId uuid.UUID, count int) (err error) {
	return err
}
