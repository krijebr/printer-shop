package repo

import (
	"context"
	"database/sql"
	"log"
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
		if filter.ProducerId != nil {
			where = " where producers.id = '" + (*filter.ProducerId).String() + "'"
		}
	}
	rows, err := p.db.QueryContext(ctx,
		"select products.id,products.name,products.price,products.status,products.created_at,producers.id,producers.name,producers.description,producers.created_at from products join producers on producers.id = producer_id"+where)
	if err != nil {
		return nil, err
	}
	products := []*entity.Product{}
	for rows.Next() {
		var product_created_at string
		var producer_created_at string
		product := new(entity.Product)
		producer := new(entity.Producer)
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Status, &product_created_at,
			&producer.Id, &producer.Name, &producer.Description, &producer_created_at)
		if err != nil {
			log.Println("Ошибка чтения строки", err)
			continue
		}
		product.CreatedAt, err = time.Parse(time.RFC3339, product_created_at)
		if err != nil {
			log.Println("Ошибка преобразования времени")
			continue
		}
		producer.CreatedAt, err = time.Parse(time.RFC3339, producer_created_at)
		if err != nil {
			log.Println("Ошибка преобразования времени")
			continue
		}
		product.Producer = producer
		products = append(products, product)
	}
	return products, nil
}
func (p *ProductRepoPg) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	row := p.db.QueryRowContext(ctx,
		"select products.id,products.name,products.price,products.status,products.created_at,producers.id,producers.name,producers.description,producers.created_at from products join producers on producers.id = producer_id where products.id=$1", id)
	var product_created_at string
	var producer_created_at string
	product := new(entity.Product)
	producer := new(entity.Producer)
	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Status, &product_created_at,
		&producer.Id, &producer.Name, &producer.Description, &producer_created_at)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrProductNotFound
		default:
			return nil, err
		}
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
	return nil
}
func (p *ProductRepoPg) DeleteById(ctx context.Context, id uuid.UUID) error {
	return nil
}
