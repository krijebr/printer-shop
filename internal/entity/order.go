package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	OrderStatus string
)

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusDone       OrderStatus = "done"
)

type Order struct {
	Id        uuid.UUID   `json:"id"`
	UserId    uuid.UUID   `json:"user_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

type OrderWithProducts struct {
	Id        uuid.UUID   `json:"id"`
	UserId    uuid.UUID   `json:"user_id"`
	Status    OrderStatus `json:"status"`
	Products  []*ProductInCart
	CreatedAt time.Time `json:"created_at"`
}
