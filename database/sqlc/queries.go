package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
	"time"
)

type userQuery interface {
	CreateUser(ctx context.Context, args CreateUserParams) (model.User, error)
	GetUserByID(ctx context.Context, id model.UserID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdatePassword(ctx context.Context, args UpdatePasswordParams) (time.Time, error)
	ListUsers(ctx context.Context, args ListUserParams) ([]model.User, error)
	DeleteUser(ctx context.Context, id model.UserID) (time.Time, error)
}

type tokenQuery interface {
	SaveRefreshToken(ctx context.Context, args SaveRefreshTokenParams) error
}

type sessionQuery interface {
	GetSession(ctx context.Context, args GetSessionsParams) (model.Session, error)
}

type roleQuery interface {
	GrantRole(ctx context.Context, args RoleParams) (model.UserRole, error)
	RevokeRole(ctx context.Context, args RoleParams) error
	GetUserRoleByID(ctx context.Context, id model.UserID) (model.UserRole, error)
	ListUsersByRole(ctx context.Context, args ListUserRoleParams) ([]model.UserRole, error)
}
type accountQuery interface {
	CreateAccount(ctx context.Context, args CreateAccountParams) (model.Account, error)
	UpdateAccount(ctx context.Context, args UpdateAccountParams) (model.Account, error)
	GetAccountByID(ctx context.Context, id model.AccountID) (model.Account, error)
	ListAccounts(ctx context.Context, args ListAccountParams) ([]model.Account, error)
	DeleteAccount(ctx context.Context, id model.AccountID) (time.Time, error)
}
type categoryQuery interface {
	CreateCategory(ctx context.Context, args CreateCategoryParams) (model.Category, error)
	UpdateCategory(ctx context.Context, args UpdateCategoryParams) (model.Category, error)
	GetCategoryByID(ctx context.Context, id model.CategoryID) (model.Category, error)
	ListCategories(ctx context.Context, args ListCategoryParams) ([]model.Category, error)
	DeleteCategory(ctx context.Context, id model.CategoryID) (time.Time, error)
}

type merchantQuery interface {
	CreateMerchant(ctx context.Context, args CreateMerchantParams) (model.Merchant, error)
	UpdateMerchant(ctx context.Context, args UpdateMerchantParams) (model.Merchant, error)
	GetMerchantByID(ctx context.Context, id model.MerchantID) (model.Merchant, error)
	ListMerchants(ctx context.Context, args ListMerchantParams) ([]model.Merchant, error)
	DeleteMerchant(ctx context.Context, id model.MerchantID) (time.Time, error)
}

type transactionQuery interface {
	CreateTransaction(ctx context.Context, args CreateTransactionParams) (model.Transaction, error)
	UpdateTransaction(ctx context.Context, args UpdateTransactionParams) (model.Transaction, error)
	GetTransactionByID(ctx context.Context, id model.TransactionID) (model.Transaction, error)
	ListTransactionsByUserID(ctx context.Context, args ListTxByUserIDParams) ([]model.Transaction, error) // we will filter with time frame
	ListTransactionsByAccountID(ctx context.Context, args ListTxByAccountIDParams) ([]model.Transaction, error)
	ListTransactionsByCategoryID(ctx context.Context, args ListTxByCategoryIDParams) ([]model.Transaction, error)
	DeleteTransaction(ctx context.Context, id model.TransactionID) (time.Time, error)
}

type QueryInterface interface {
	userQuery
	tokenQuery
	sessionQuery
	roleQuery
	accountQuery
	categoryQuery
	merchantQuery
	transactionQuery
}

// we want to ensure all our methods in the interface are implemented by our Queries struct
var _ QueryInterface = (*Queries)(nil)
