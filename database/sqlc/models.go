package database

import "time"

// UserID is identifier for our User
type UserID string

// User struct contains fields for our users which represents our user object
type User struct {
	ID UserID `json:"id" db:"user_id"`
	Email string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	PasswordChangedAt time.Time `json:"password_changed_at" db:"password_changed_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}
