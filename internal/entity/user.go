package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserStatus string
	UserRole   string
)

const (
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"

	UserRoleGuest    UserRole = "guest"
	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
)

type User struct {
	Id           uuid.UUID  `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Status       UserStatus `json:"status"`
	Role         UserRole   `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
}
type UserFilter struct {
	UserStatus *UserStatus `json:"user_status"`
	UserRole   *UserRole   `json:"user_role"`
}
