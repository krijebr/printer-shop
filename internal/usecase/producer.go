package usecase

import (
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)

type producer struct {
}

func NewProducer() Producer {
	return &producer{}
}
func (p *producer) GetAll() ([]*entity.Producer, error) {
	return nil, err
}
func (p *producer) GetById(id uuid.UUID) (*entity.Producer, error) {
	return nil, err
}
func (p *producer) Create(producer entity.Producer) (*entity.Producer, error) {
	return nil, err
}
func (p *producer) UpdateById(id uuid.UUID, producer entity.Producer) (*entity.Producer, error) {
	return nil, err
}
func (p *producer) DeleteById(id uuid.UUID) error {
	return err
}
