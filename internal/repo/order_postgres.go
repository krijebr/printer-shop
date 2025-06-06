package repo

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	_ "github.com/lib/pq"
)

type OrderRepoPg struct {
	db *sql.DB
}

func NewOrderRepoPg(db *sql.DB) Order {
	return &OrderRepoPg{
		db: db,
	}
}

func (o *OrderRepoPg) Create(ctx context.Context, order *entity.Order) error {
	_, err := o.db.ExecContext(ctx, "insert into orders(id, user_id, status, created_at) values ($1,$2,$3,$4)",
		order.Id, order.UserId, order.Status, order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (o *OrderRepoPg) GetAll(ctx context.Context, filter *entity.OrderFilter) ([]*entity.Order, error) {
	var orderCreatedAt string
	var productCreatedAt string
	var producerCreatedAt string
	where := ""
	if filter != nil {
		whereS := []string{}
		if filter.UserId != nil {
			whereS = append(whereS, "orders.user_id = '"+(*filter.UserId).String()+"'")
		}
		if filter.Status != nil {
			whereS = append(whereS, "orders.status = '"+string(*filter.Status)+"'")
		}
		where = " where " + strings.Join(whereS, " and ")
	}
	rows, err := o.db.QueryContext(ctx, "select orders.id, orders.user_id, orders.status, orders.created_at, products.id, products.name, order_products.product_price, producers.id, producers.name, producers.description, producers.created_at, products.status, products.created_at, order_products.product_count from orders join order_products on order_products.order_id = orders.id join products on order_products.product_id = products.id join producers on products.producer_id = producers.id"+where)
	if err != nil {
		return nil, err
	}
	orders := []*entity.Order{}
	previousOrderId := uuid.Nil
	for rows.Next() {
		order := entity.Order{}
		product := &entity.ProductInCart{
			Product: &entity.Product{},
		}
		producer := new(entity.Producer)
		err := rows.Scan(&order.Id, &order.UserId, &order.Status, &orderCreatedAt, &product.Product.Id, &product.Product.Name,
			&product.Product.Price, &producer.Id, &producer.Name, &producer.Description, &producerCreatedAt,
			&product.Product.Status, &productCreatedAt, &product.Count)
		if err != nil {
			return nil, err
		}
		if order.Id != previousOrderId {
			order.CreatedAt, err = time.Parse(time.RFC3339, orderCreatedAt)
			if err != nil {
				return nil, err
			}
			slog.Debug(string(order.Status) + string(order.Status))

			orders = append(orders, &order)
			previousOrderId = order.Id
		}
		product.Product.CreatedAt, err = time.Parse(time.RFC3339, productCreatedAt)
		if err != nil {
			return nil, err
		}
		producer.CreatedAt, err = time.Parse(time.RFC3339, producerCreatedAt)
		if err != nil {
			return nil, err
		}
		product.Product.Producer = producer
		orders[len(orders)-1].Products = append(orders[len(orders)-1].Products, product)
	}
	return orders, nil
}
func (o *OrderRepoPg) GetById(ctx context.Context, id uuid.UUID) (order *entity.Order, err error) {
	return nil, nil
}
func (o *OrderRepoPg) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	return nil
}
func (o *OrderRepoPg) UpdateById(ctx context.Context, id uuid.UUID) (err error) {
	return nil
}
