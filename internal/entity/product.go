package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	ProductStatusPublished ProductStatus = "published"
	ProductStatusHidden    ProductStatus = "hidden"
)

type (
	ProductStatus string

	Product struct {
		Id        uuid.UUID     `json:"id"`
		Name      string        `json:"name"`
		Price     float32       `jsone:"price"`
		Producer  *Producer     `json:"producer"`
		Status    ProductStatus `json:"status"`
		CreatedAt time.Time     `json:"created_at"`
	}

	ProductInCart struct {
		Product Product `json:"product"`
		Count   int     `json:"count"`
	}

	ProductFilter struct {
		ProducerId *uuid.UUID `json:"producer_id"`
	}
)
