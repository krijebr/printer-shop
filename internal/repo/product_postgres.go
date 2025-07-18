package repo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	_ "github.com/lib/pq"
)

type ProductRepoPg struct {
	db *sql.DB
}

func NewProductRepoPg(db *sql.DB) Product {
	return &ProductRepoPg{
		db: db,
	}
}

func (p *ProductRepoPg) GetAll(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, error) {

	where := ""
	if filter != nil {
		whereS := []string{}
		if filter.ProducerId != nil {
			whereS = append(whereS, "producers.id = '"+(*filter.ProducerId).String()+"'")
		}
		if filter.Status != nil {
			whereS = append(whereS, "products.status = '"+string(*filter.Status)+"'")
		}
		where = " where " + strings.Join(whereS, " and ")
	}

	rows, err := p.db.QueryContext(ctx,
		"select "+
			"products.id,products.name,products.price,products.status,products.created_at,producers.id,producers.name,producers.description,producers.created_at "+
			"from "+
			"products join producers on producers.id = producer_id"+
			where)
	if err != nil {
		return nil, err
	}
	products := []*entity.Product{}
	for rows.Next() {
		product, err := p.scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepoPg) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	row := p.db.QueryRowContext(ctx,
		"select "+
			"products.id,products.name,products.price,products.status,products.created_at,producers.id,producers.name,producers.description,producers.created_at "+
			"from "+
			"products join producers on producers.id = producer_id "+
			"where products.id=$1",
		id)
	product, err := p.scanProduct(row)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrProductNotFound
		default:
			return nil, err
		}
	}
	return product, nil
}

func (p *ProductRepoPg) Create(ctx context.Context, product entity.Product) error {
	_, err := p.db.ExecContext(ctx, "insert into products (id, name, price, producer_id, status, created_at) values ($1,$2,$3,$4,$5,$6)",
		product.Id, product.Name, product.Price, product.Producer.Id, product.Status, product.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepoPg) Update(ctx context.Context, product entity.Product) error {
	set := []string{}
	if product.Name != "" {
		set = append(set, "name = '"+product.Name+"'")
	}
	if product.Price != 0 {
		set = append(set, "price = "+strconv.FormatFloat(float64(product.Price), 'f', 2, 32))
	}
	if product.Producer.Id != uuid.Nil {
		set = append(set, "producer_id = '"+product.Producer.Id.String()+"'")
	}
	if product.Status != "" {
		set = append(set, "status = '"+string(product.Status)+"'")
	}
	_, err := p.db.ExecContext(ctx, "update products set "+strings.Join(set, ", ")+" where id = $1", product.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepoPg) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, "delete from products where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepoPg) scanProduct(row Row) (*entity.Product, error) {
	var productCreatedAt string
	var producerCreatedAt string
	product := new(entity.Product)
	producer := new(entity.Producer)
	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Status, &productCreatedAt,
		&producer.Id, &producer.Name, &producer.Description, &producerCreatedAt)
	if err != nil {
		return nil, err
	}
	product.CreatedAt, err = time.Parse(time.RFC3339, productCreatedAt)
	if err != nil {
		return nil, err
	}
	producer.CreatedAt, err = time.Parse(time.RFC3339, producerCreatedAt)
	if err != nil {
		return nil, err
	}
	product.Producer = producer
	return product, nil
}
