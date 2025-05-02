package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	ProductStatus string
	UserStatus    string
	OrderStatus   string
	UserRole      string
)

const (
	ProductStatusPublished ProductStatus = "published"
	ProductStatusHidden    ProductStatus = "hidden"

	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"

	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusDone       OrderStatus = "done"

	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
)

type User struct {
	Id        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Status    UserStatus `json:"status"`
	Role      UserRole   `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
}

type Product struct {
	Id         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Price      float32       `jsone:"price"`
	ProducerId uuid.UUID     `json:"producer_id"`
	Status     ProductStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
}

type Producer struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Cart struct {
	UserId    uuid.UUID `json:"user_id"`
	ProductId uuid.UUID `json:"product_id"`
	Count     int       `json:"count"`
}

type Order struct {
	Id        uuid.UUID   `json:"id"`
	UserId    uuid.UUID   `json:"user_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

type OrderProduct struct {
	ProductId    uuid.UUID `json:"product_id"`
	OrderId      uuid.UUID `json:"order_id"`
	ProductCount int       `json:"product_count"`
	ProductPrice float32   `json:"product_price"`
}
