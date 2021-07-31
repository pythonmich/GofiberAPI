package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
)

const grantRole = `--name: GrantRole :exec
INSERT INTO user_roles (user_id, role)
VALUES ($1, $2)
RETURNING user_id, role, created_at
`

type RoleParams struct {
	UserID model.UserID `json:"user_id"`
	Role   model.Role   `json:"role"`
}

func (q *Queries) GrantRole(ctx context.Context, args RoleParams) (model.UserRole, error) {
	q.logs.WithField("func", "database/sqlc/user_roles.go -> GrantRole()").Debug()
	row := q.db.QueryRowContext(ctx, grantRole, args.UserID, args.Role)
	var userRole model.UserRole
	err := row.Scan(
		&userRole.UserID,
		&userRole.Role,
		&userRole.CreatedAt,
	)
	if err != nil {
		q.logs.WithError(err).Warn()
		return model.UserRole{}, err
	}
	return userRole, nil
}

const revokeRole = `--name: RevokeRole :exec
DELETE FROM user_roles
WHERE user_id = $1
AND role = $2
`

func (q *Queries) RevokeRole(ctx context.Context, args RoleParams) error {
	q.logs.WithField("func", "database/sqlc/user_roles.go -> RevokeRole()").Debug()
	_, err := q.GetUserRoleByID(ctx, args.UserID)
	if err != nil {
		return err
	}
	_, err = q.db.ExecContext(ctx, revokeRole, args.UserID, args.Role)
	return err
}

const getUserRoleByID = `--name: GetUserRoleByID :one
SELECT role, created_at FROM user_roles
WHERE user_id = $1
LIMIT 1
`

func (q Queries) GetUserRoleByID(ctx context.Context, id model.UserID) (model.UserRole, error) {
	q.logs.WithField("func", "database/sqlc/user_roles.go -> GetUserRoleByID()").Debug()
	row := q.db.QueryRowContext(ctx, getUserRoleByID, id)
	var userRole model.UserRole
	err := row.Scan(
		&userRole.Role,
		&userRole.CreatedAt,
	)
	if err != nil {
		q.logs.WithError(err).Warn()
		return model.UserRole{}, err
	}
	return userRole, nil
}

type ListUserRoleParams struct {
	UserID model.UserID `json:"user_id"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
}

const listUserByRole = `--name: ListUsersByRole :many
SELECT * FROM user_roles
WHERE user_id = $1
LIMIT $2
OFFSET $3
`

func (q *Queries) ListUsersByRole(ctx context.Context, args ListUserRoleParams) ([]model.UserRole, error) {
	q.logs.WithField("func", "database/sqlc/user_roles.go -> ListUsersByRole()").Debug()

	rows, err := q.db.QueryContext(ctx, listUserByRole, args.UserID, args.Limit, args.Offset)
	if err != nil {
		q.logs.WithError(err).Warn()
		return nil, err
	}
	defer func() {
		q.logs.WithField("func", "database/sqlc/user_roles.go -> ListUsersByRole()->func()").Debug()
		// row error
		err = rows.Close()
		if err != nil {
			q.logs.WithError(err).Warn()
		}
		q.logs.Debug("rows closed successfully")
	}()
	var userRoles []model.UserRole
	for rows.Next() {
		var userRole model.UserRole
		err = rows.Scan(
			&userRole.UserID,
			&userRole.Role,
			&userRole.CreatedAt,
		)
		userRoles = append(userRoles, userRole)
	}
	q.logs.Debug("successful")
	return userRoles, err
}
