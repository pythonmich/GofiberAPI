package models

import (
	"FiberFinanceAPI/utils"
	"time"
)

// AccountID is our identifier for our account
type AccountID string

// AccountType its the type of Account
type AccountType string

const (
	Cash   AccountType = "cash"
	Credit AccountType = "credit"
)

// Account represents our user account model or structure
type Account struct {
	AccountID AccountID          `json:"account_id"`
	UserID    UserID             `json:"user_id"`
	Name      string             `json:"account_name"`
	Type      AccountType        `json:"account_type"`
	Balance   int64              `json:"balance"`
	Currency  utils.CurrencyCode `json:"currency"`
	CreatedAt time.Time          `json:"created_at"`
	DeletedAt time.Time          `json:"-"`
}
