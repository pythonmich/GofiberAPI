package database

import (
	"context"
)

type QueryInterface interface {
	CreateUser(ctx context.Context, args CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	SaveRefreshToken(ctx context.Context, args SaveRefreshTokenParams) error
	GetSession(ctx context.Context, args GetSessionsParams) (Session, error)
	GrantRole(ctx context.Context, args RoleParams) (UserRole, error)
	RevokeRole(ctx context.Context, args RoleParams) error
	GetUserRoleByID(ctx context.Context, id UserID) (UserRole, error)
	ListUsersByRole(ctx context.Context, args ListUserRoleParams) ([]UserRole, error)
}

var _ QueryInterface = (*Queries)(nil)
