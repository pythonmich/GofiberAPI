package models

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
