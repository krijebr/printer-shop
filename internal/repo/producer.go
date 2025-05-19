package repo

import (
	"context"
	"database/sql"

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

func (p *ProducerRepoPg) GetAll(ctx context.Context) (allProducers []*entity.Producer, err error) {
	return nil, nil
}
func (p *ProducerRepoPg) GetById(ctx context.Context, id uuid.UUID) (producer *entity.Producer, err error) {
	return nil, nil
}
func (p *ProducerRepoPg) Create(ctx context.Context, producer entity.Producer) (createdProducer *entity.Producer, err error) {
	return nil, nil
}
func (p *ProducerRepoPg) Update(ctx context.Context, producer entity.Producer) (updatedProducer *entity.Producer, err error) {
	return nil, nil
}
func (p *ProducerRepoPg) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	return nil
}
