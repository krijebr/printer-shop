package entity

import (
	"github.com/google/uuid"
)

type Cart struct {
	UserId    uuid.UUID `json:"user_id"`
	ProductId uuid.UUID `json:"product_id"`
	Count     int       `json:"count"`
}
