package database

import (
	"context"
)
type QueryInterface interface {
	CreateUser(ctx context.Context, args CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

var _ QueryInterface = (*Queries)(nil)