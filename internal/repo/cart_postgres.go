package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	_ "github.com/lib/pq"
)

type CartRepoPg struct {
	db *sql.DB
}

func NewCartRepoPg(db *sql.DB) Cart {
	return &CartRepoPg{
		db: db,
	}
}

func (c *CartRepoPg) GetAllProducts(ctx context.Context, userId uuid.UUID) ([]*entity.ProductInCart, error) {
	rows, err := c.db.QueryContext(ctx,
		"select products.id,products.name,products.price,products.status,products.created_at,producers.id,producers.name,producers.description,producers.created_at, carts.count from carts join products on carts.product_id = products.id join producers on producers.id = producer_id where carts.user_id=$1", userId)
	if err != nil {
		return nil, err
	}
	productsInCart := []*entity.ProductInCart{}
	for rows.Next() {
		productInCart, err := c.scanProductInCart(rows)
		if err != nil {
			return nil, err
		}
		productsInCart = append(productsInCart, productInCart)
	}
	return productsInCart, nil
}
func (c *CartRepoPg) AddProduct(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) error {
	return nil
}
func (c *CartRepoPg) UpdateCount(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) error {
	return nil
}
func (c *CartRepoPg) scanProductInCart(row Row) (*entity.ProductInCart, error) {
	var product_created_at string
	var producer_created_at string
	productInCart := new(entity.ProductInCart)
	product := new(entity.Product)
	producer := new(entity.Producer)
	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Status, &product_created_at, &producer.Id, &producer.Name, &producer.Description, &producer_created_at, &productInCart.Count)
	if err != nil {
		return nil, err
	}
	product.CreatedAt, err = time.Parse(time.RFC3339, product_created_at)
	if err != nil {
		return nil, err
	}
	producer.CreatedAt, err = time.Parse(time.RFC3339, producer_created_at)
	if err != nil {
		return nil, err
	}
	product.Producer = producer
	productInCart.Product = product
	return productInCart, nil
}
