package repo

import (
	"context"
	"database/sql"
	"errors"
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
	_, err := c.db.ExecContext(ctx, "insert into carts (user_id, product_id, count) values ($1,$2,$3)", userId, productId, count)
	if err != nil {
		return err
	}
	return nil
}
func (c *CartRepoPg) UpdateCount(ctx context.Context, userId uuid.UUID, productId uuid.UUID, count int) error {
	_, err := c.db.ExecContext(ctx, "update carts set count = $1 where user_id = $2 and product_id = $3", count, userId, productId)
	if err != nil {
		return err
	}
	return nil
}

func (c *CartRepoPg) DeleteByProductId(ctx context.Context, userId uuid.UUID, productId uuid.UUID) error {
	_, err := c.db.ExecContext(ctx, "delete from carts where user_id = $1 and product_id = $2", userId, productId)
	if err != nil {
		return err
	}
	return nil
}
func (c *CartRepoPg) GetProductCountById(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (int, error) {
	row := c.db.QueryRowContext(ctx, "select count from carts where user_id=$1 and product_id = $2", userId, productId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, nil
		default:
			return 0, err
		}
	}
	return count, nil
}
func (c *CartRepoPg) CheckIfExistsById(ctx context.Context, productId uuid.UUID) (bool, error) {
	row := c.db.QueryRowContext(ctx, "select count(count) from carts where product_id = $1", productId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
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
func (c *CartRepoPg) ClearCart(ctx context.Context, userId uuid.UUID) error {
	_, err := c.db.ExecContext(ctx, "delete from carts where user_id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}
