package models

import "time"

// Role is a function can serve
type Role string

const (
	// RoleAdmin is the administrator of our app
	RoleAdmin Role = "admin"
)

type UserRole struct {
	UserID    UserID    `json:"-"`
	Role      Role      `json:"role" validate:"required,alphanum"`
	CreatedAt time.Time `json:"created_at"`
}
