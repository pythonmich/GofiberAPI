package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
	"time"
)

const createMerchant = `--name: CreateMerchant :one
INSERT INTO merchant(user_id, name) 
VALUES($1, $2)
RETURNING merchant_id, user_id, name, created_at, deleted_at`

type CreateMerchantParams struct {
	UserID model.UserID `json:"user_id"`
	Name   string       `json:"name"`
}

func (q *Queries) CreateMerchant(ctx context.Context, args CreateMerchantParams) (model.Merchant, error) {
	q.logs.WithField("func", "database/sqlc/merchant.go -> CreateMerchant()").Debug()
	row := q.db.QueryRowContext(ctx, createMerchant, args.UserID, args.Name)
	var merchant model.Merchant
	err := row.Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.Name,
		&merchant.CreatedAt,
		&merchant.DeletedAt,
	)
	return merchant, err
}

const updateMerchant = `--name: UpdateMerchant :one
UPDATE merchant SET name = $2
WHERE merchant_id = $1 
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING merchant_id, user_id, name, created_at, deleted_at
`

type UpdateMerchantParams struct {
	MerchantID model.MerchantID `json:"merchant_id"`
	Name       string           `json:"name"`
}

func (q *Queries) UpdateMerchant(ctx context.Context, args UpdateMerchantParams) (model.Merchant, error) {
	q.logs.WithField("func", "database/sqlc/merchant.go -> UpdateMerchant()").Debug()
	row := q.db.QueryRowContext(ctx, updateMerchant, args.MerchantID, args.Name)
	var merchant model.Merchant
	err := row.Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.Name,
		&merchant.CreatedAt,
		&merchant.DeletedAt,
	)
	return merchant, err
}

const getMerchant = `--name: GetMerchantByID :one
SELECT * FROM merchant
WHERE merchant_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1`

func (q *Queries) GetMerchantByID(ctx context.Context, id model.MerchantID) (model.Merchant, error) {
	q.logs.WithField("func", "database/sqlc/merchant.go -> GetMerchantByID()").Debug()
	row := q.db.QueryRowContext(ctx, getMerchant, id)
	var merchant model.Merchant
	err := row.Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.Name,
		&merchant.CreatedAt,
		&merchant.DeletedAt,
	)
	return merchant, err
}

const listMerchants = `--name: ListMerchants :one
SELECT * FROM merchant
WHERE user_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY merchant_id
LIMIT $2
OFFSET $3`

type ListMerchantParams struct {
	UserID model.UserID `json:"user_id"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
}

func (q *Queries) ListMerchants(ctx context.Context, args ListMerchantParams) ([]model.Merchant, error) {
	q.logs.WithField("func", "database/sqlc/merchant.go -> ListMerchants()").Debug()
	rows, err := q.db.QueryContext(ctx, listMerchants, args.UserID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		q.logs.WithError(err).Warn()
	}()
	var merchants []model.Merchant
	for rows.Next() {
		var merchant model.Merchant
		err = rows.Scan(
			&merchant.ID,
			&merchant.UserID,
			&merchant.Name,
			&merchant.CreatedAt,
			&merchant.DeletedAt,
		)
		merchants = append(merchants, merchant)
	}
	return merchants, err
}

const deleteMerchant = `--name: DeleteMerchant :one
UPDATE merchant SET deleted_at = now()
WHERE merchant_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;`

func (q *Queries) DeleteMerchant(ctx context.Context, id model.MerchantID) (time.Time, error) {
	q.logs.WithField("func", "database/sqlc/merchant.go -> DeleteMerchant()").Debug()
	row := q.db.QueryRowContext(ctx, deleteMerchant, id)
	var merchant model.Merchant
	err := row.Scan(
		&merchant.DeletedAt,
	)
	return merchant.DeletedAt, err
}
