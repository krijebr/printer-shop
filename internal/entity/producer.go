package entity

import (
	"time"

	"github.com/google/uuid"
)

type Producer struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
