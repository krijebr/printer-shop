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
	UserId    *User       `json:"omitempty"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Products  []*ProductInCart
}
type OrderFilter struct {
	UserId *uuid.UUID
	Status *OrderStatus
}
