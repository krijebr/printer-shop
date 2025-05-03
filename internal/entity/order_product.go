package entity

import (
	"github.com/google/uuid"
)

type OrderProduct struct {
	ProductId    uuid.UUID `json:"product_id"`
	OrderId      uuid.UUID `json:"order_id"`
	ProductCount int       `json:"product_count"`
	ProductPrice float32   `json:"product_price"`
}
