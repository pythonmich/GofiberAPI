package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
	"time"
)

const createUser = `--name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING user_id, email, password_hash, password_changed_at, created_at, deleted_at  
`

type CreateUserParams struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, args CreateUserParams) (model.User, error) {
	q.logs.WithField("func", "database/sqlc/users.go -> CreateUser()").Debug()
	row := q.db.QueryRowContext(ctx, createUser, args.Email, args.PasswordHash)
	var user model.User
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
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1`

func (q *Queries) GetUserByID(ctx context.Context, id model.UserID) (model.User, error) {
	q.logs.WithField("func", "database/sqlc/users.go -> GetUserByID()").Debug()
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var user model.User
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
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	q.logs.WithField("func", "database/sqlc/users.go -> GetUserByEmail()").Debug(email)
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var user model.User
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

const listUsers = `--name: ListUsers :many
SELECT * FROM users
WHERE deleted_at = '0001-01-01 00:00:00Z'
ORDER BY user_id
LIMIT $1
OFFSET $2
`

type ListUserParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, args ListUserParams) ([]model.User, error) {
	q.logs.WithField("func", "database/sqlc/user.go -> ListUsers()").Debug()
	rows, err := q.db.QueryContext(ctx, listUsers, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		q.logs.WithError(err).Warn()
	}()
	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.PasswordChangedAt,
			&user.CreatedAt,
			&user.DeletedAt,
		)
		users = append(users, user)
	}
	return users, err
}

type UpdatePasswordParams struct {
	UserID       model.UserID `json:"user_id"`
	HashPassword string       `json:"hash_password"`
}

const updatePassword = `--name: UpdatePassword :one
UPDATE users SET password_changed_at = now(),
password_hash = $2
WHERE user_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING password_changed_at`

func (q Queries) UpdatePassword(ctx context.Context, args UpdatePasswordParams) (time.Time, error) {
	q.logs.WithField("func", "database/sqlc/user.go -> UpdatePassword()").Debug()
	row := q.db.QueryRowContext(ctx, updatePassword, args.UserID, args.HashPassword)
	var user model.User
	err := row.Scan(
		&user.PasswordChangedAt,
	)
	return user.PasswordChangedAt, err
}

const deleteUser = `--name: DeleteUser :exec
UPDATE users SET deleted_at = now(),
email = concat(email, '-DELETED-', uuid_generate_v4())
WHERE user_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;
`

func (q *Queries) DeleteUser(ctx context.Context, id model.UserID) (time.Time, error) {
	q.logs.WithField("func", "database/sqlc/user.go -> DeleteUser()").Debug()
	row := q.db.QueryRowContext(ctx, deleteUser, id)
	var user model.User
	err := row.Scan(
		&user.DeletedAt,
	)
	return user.DeletedAt, err
}
