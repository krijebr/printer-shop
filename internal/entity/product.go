package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	ProductStatus string
)

const (
	ProductStatusPublished ProductStatus = "published"
	ProductStatusHidden    ProductStatus = "hidden"
)

type Product struct {
	Id        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Price     float32       `jsone:"price"`
	Producer  *Producer     `json:"omitempty"`
	Status    ProductStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

type ProductInCart struct {
	Product Product
	Count   int `json:"count"`
}
type ProductForOrder struct {
	ProductId uuid.UUID
	Count     int
}
type ProductFilter struct {
	ProductId *uuid.UUID
}
