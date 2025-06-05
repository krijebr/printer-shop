package repo

import (
	"context"
	"database/sql"
	"strings"

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

func (o *OrderRepoPg) Create(ctx context.Context, userId uuid.UUID, products []*entity.ProductInCart) (err error) {
	return nil
}
func (o *OrderRepoPg) GetAll(ctx context.Context, filter *entity.OrderFilter) (allOrders []*entity.Order, err error) {
	where := ""
	if filter != nil {
		whereS := []string{}
		if filter.UserId != nil {
			whereS = append(whereS, "status = '"+(*filter.UserId).String()+"'")
		}
		if filter.Status != nil {
			whereS = append(whereS, "role = '"+string(*filter.Status)+"'")
		}
		where = " where " + strings.Join(whereS, " and ")
	}
	rows, err := o.db.QueryContext(ctx, "select * from orders"+where)
	if err != nil {
		return nil, err
	}
	orders := []*entity.Order{}
	for rows.Next() {
		order, err := o.scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
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

func (o *OrderRepoPg) scanOrder(row Row) (*entity.Order, error) {

	return nil, nil
}
