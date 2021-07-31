package models

import "time"

// TransactionID is our identifier for our transactions
type TransactionID string

// TransactionType its the type of Transaction
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID              TransactionID   `json:"id"`
	UserID          UserID          `json:"user_id"`
	AccountID       AccountID       `json:"account_id"`
	CategoryID      CategoryID      `json:"category_id"`
	Name            string          `json:"name"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          int64           `json:"amount"`
	Notes           string          `json:"notes"`
	Date            time.Time       `json:"date"`
	CreatedAt       time.Time       `json:"created_at"`
	DeletedAt       time.Time       `json:"-"`
}
