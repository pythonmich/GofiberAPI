package database

import (
	"context"
	logs "github.com/sirupsen/logrus"
)

const createUser = `--name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING user_id, email, password_hash, password_changed_at, created_at, deleted_at  
`

type CreateUserParams struct {
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, args CreateUserParams) (User, error) {
	logs.WithField("func", "database/sqlc/users.go -> CreateUser()").Debug()
	row := q.db.QueryRowContext(ctx, createUser, args.Email, args.PasswordHash)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.PasswordChangedAt,
		&user.CreatedAt,
		&user.DeletedAt,
		)
	return user, err
}

const getUserByID = `--name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1
LIMIT 1`

func (q *Queries) GetUserByID(ctx context.Context, userID string) (User, error) {
	logs.WithField("func", "database/sqlc/users.go -> GetUserByID()").Debug()
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.PasswordChangedAt,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	return user, err
}

const getUserByEmail = `--name: GetUserByID :one
SELECT * FROM users
WHERE email = $1
LIMIT 1`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	logs.WithField("func", "database/sqlc/users.go -> GetUserByEmail()").Debug()
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.PasswordChangedAt,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	return user, err
}