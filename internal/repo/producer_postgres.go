package repo

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	_ "github.com/lib/pq"
)

type ProducerRepoPg struct {
	db *sql.DB
}

func NewProducerRepoPg(db *sql.DB) Producer {
	return &ProducerRepoPg{
		db: db,
	}
}

func (p *ProducerRepoPg) GetAll(ctx context.Context) ([]*entity.Producer, error) {
	rows, err := p.db.QueryContext(ctx, "select * from producers")
	if err != nil {
		return nil, err
	}
	producers := []*entity.Producer{}
	for rows.Next() {
		producer, err := p.scanProducer(rows)
		if err != nil {
			return nil, err
		}
		producers = append(producers, producer)
	}
	return producers, nil
}
func (p *ProducerRepoPg) GetById(ctx context.Context, id uuid.UUID) (*entity.Producer, error) {
	row := p.db.QueryRowContext(ctx, "select * from producers where id = $1", id)
	producer, err := p.scanProducer(row)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrProducerNotFound
		default:
			return nil, err
		}
	}
	return producer, nil
}
func (p *ProducerRepoPg) Create(ctx context.Context, producer entity.Producer) error {
	_, err := p.db.ExecContext(ctx, "insert into producers (id, name, description, created_at) values ($1,$2,$3,$4)",
		producer.Id, producer.Name, producer.Description, producer.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (p *ProducerRepoPg) Update(ctx context.Context, producer entity.Producer) error {
	set := []string{}
	if producer.Name != "" {
		set = append(set, "name = '"+producer.Name+"'")
	}
	if producer.Description != "" {
		set = append(set, "description = '"+producer.Description+"'")
	}

	_, err := p.db.ExecContext(ctx, "update producers set "+strings.Join(set, ", ")+" where id = $1", producer.Id)
	if err != nil {
		return err
	}
	return nil
}
func (p *ProducerRepoPg) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, "delete from producers where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProducerRepoPg) scanProducer(row Row) (*entity.Producer, error) {
	var dateStr string
	producer := new(entity.Producer)
	err := row.Scan(&producer.Id, &producer.Name, &producer.Description, &dateStr)
	if err != nil {
		return nil, err
	}
	producer.CreatedAt, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
