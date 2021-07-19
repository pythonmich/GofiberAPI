package database

import "time"

// UserID is identifier for our User
type UserID string

// User struct contains fields for our users which represents our user object
type User struct {
	ID                UserID    `json:"id" db:"user_id"`
	Email             string    `json:"email"`
	PasswordHash      string    `json:"-"`
	PasswordChangedAt time.Time `json:"-"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	DeletedAt         time.Time `json:"-"`
}

type DeviceID string

// Session represents our User's session
type Session struct {
	UserID       UserID    `json:"user_id"`
	DeviceID     DeviceID  `json:"device_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    int64     `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// LoginSession contains our device id
type LoginSession struct {
	DeviceID DeviceID `json:"device_id" validate:"required"`
}

// Role is a function can serve
type Role string

const (
	// RoleAdmin is the administrator of our app
	RoleAdmin Role = "admin"
)

type UserRole struct {
	UserID    UserID    `json:"-"`
	Role      Role      `json:"role" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
