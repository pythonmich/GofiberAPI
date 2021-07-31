package models

import "time"

// MerchantID is our identifier for our merchant
type MerchantID string

// Merchant represents our user merchant model or structure
type Merchant struct {
	ID        MerchantID `json:"id"`
	UserID    UserID     `json:"user_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt time.Time  `json:"-"`
}
