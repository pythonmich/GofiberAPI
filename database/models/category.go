package models

import "time"

// CategoryID is our identifier for our category
type CategoryID string

// Category represents our user category model or structure
type Category struct {
	ID        CategoryID `json:"id"`
	ParentID  CategoryID `json:"parent_id"`
	UserID    UserID     `json:"user_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt time.Time  `json:"-"`
}
