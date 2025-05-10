package usecase

import (
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type cart struct {
}

func NewCart() Cart {
	return &cart{}
}
func (c *cart) GetAllProducts() (allProducts []*entity.ProductInCart, err error) {
	return nil, err
}
func (c *cart) AddProduct(productId uuid.UUID, count int) (err error) {
	return err
}
func (c *cart) UpdateCount(productId uuid.UUID, count int) (err error) {
	return err
}
