package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusDone       OrderStatus = "done"
)

type (
	OrderStatus string

	Order struct {
		Id        uuid.UUID   `json:"id"`
		UserId    *User       `json:"omitempty"`
		Status    OrderStatus `json:"status"`
		CreatedAt time.Time   `json:"created_at"`
		Products  []*ProductInOrder
	}
	OrderFilter struct {
		UserId *uuid.UUID   `json:"user_id"`
		Status *OrderStatus `json:"order_status"`
	}
)
