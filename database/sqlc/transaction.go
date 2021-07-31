package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
	"time"
)

const createTransaction = `--name: CreateTransaction :one
INSERT INTO transactions (user_id, account_id, category_id, name, transaction_type, amount, notes, date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING transaction_id, user_id, account_id, category_id, name, transaction_type, amount, notes, date, created_at, deleted_at`

type CreateTransactionParams struct {
	UserID          model.UserID          `json:"user_id"`
	AccountID       model.AccountID       `json:"account_id"`
	CategoryID      model.CategoryID      `json:"category_id"`
	Name            string                `json:"name"`
	TransactionType model.TransactionType `json:"transaction_type"`
	Amount          int64                 `json:"amount"`
	Notes           string                `json:"notes"`
	Date            time.Time             `json:"date"`
}

func (q *Queries) CreateTransaction(ctx context.Context, args CreateTransactionParams) (model.Transaction, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> CreateTransaction()").Debug()

	row := q.db.QueryRowContext(ctx, createTransaction, args.UserID, args.AccountID, args.CategoryID, args.Name,
		args.TransactionType, args.Amount, args.Notes, args.Date)
	var transaction model.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AccountID,
		&transaction.CategoryID,
		&transaction.Name,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.Notes,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.DeletedAt,
	)
	return transaction, err
}

const updateTransaction = `--name: UpdateTransaction :one
UPDATE transactions SET account_id = $2, 
category_id = $3, 
name = $4, 
transaction_type = $5,  
amount = $6,
notes = $7,
date = $8
WHERE transaction_id = $1
RETURNING transaction_id, user_id, account_id, category_id, name, transaction_type, amount, notes, date, created_at, deleted_at`

type UpdateTransactionParams struct {
	TransactionID   model.TransactionID   `json:"transaction_id"`
	UserID          model.UserID          `json:"user_id"`
	AccountID       model.AccountID       `json:"account_id"`
	CategoryID      model.CategoryID      `json:"category_id"`
	Name            string                `json:"name"`
	TransactionType model.TransactionType `json:"transaction_type"`
	Amount          int64                 `json:"amount"`
	Notes           string                `json:"notes"`
	Date            time.Time             `json:"date"`
}

func (q *Queries) UpdateTransaction(ctx context.Context, args UpdateTransactionParams) (model.Transaction, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> UpdateTransaction()").Debug()
	row := q.db.QueryRowContext(ctx, updateTransaction, args.TransactionID, args.UserID, args.AccountID, args.CategoryID, args.Name,
		args.TransactionType, args.Amount, args.Notes, args.Date)
	var transaction model.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AccountID,
		&transaction.CategoryID,
		&transaction.Name,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.Notes,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.DeletedAt,
	)
	return transaction, err
}

const getTransaction = `--name: GetTransactionByID :one
SELECT * FROM transactions
WHERE transaction_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1`

func (q *Queries) GetTransactionByID(ctx context.Context, id model.TransactionID) (model.Transaction, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> GetTransactionByID()").Debug()
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var transaction model.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AccountID,
		&transaction.CategoryID,
		&transaction.Name,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.Notes,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.DeletedAt,
	)
	return transaction, err

}

const listTXByUserID = `--name: ListTransactionsByUserID :many
SELECT * FROM transactions
WHERE user_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
AND date > $4
AND date < $5
ORDER BY transaction_id
LIMIT  $2
OFFSET $3`

type ListTxByUserIDParams struct {
	UserID model.UserID `json:"user_id"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
	From   time.Time    `json:"from"`
	To     time.Time    `json:"to"`
}

func (q *Queries) ListTransactionsByUserID(ctx context.Context, args ListTxByUserIDParams) ([]model.Transaction, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> ListTransactionsByUserID()").Debug()
	rows, err := q.db.QueryContext(ctx, listTXByUserID, args.UserID, args.Limit, args.Offset, args.From, args.To)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		q.logs.WithError(err).Warn()
	}()
	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.CategoryID,
			&transaction.Name,
			&transaction.TransactionType,
			&transaction.Amount,
			&transaction.Notes,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.DeletedAt,
		)
		transactions = append(transactions, transaction)
	}
	return transactions, err
}

const listTXByAccountID = `--name: ListTransactionsByAccountID :many
SELECT * FROM transactions
WHERE account_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
AND date > $4
AND date < $5
ORDER BY transaction_id
LIMIT  $2
OFFSET $3`

type ListTxByAccountIDParams struct {
	AccountID model.AccountID `json:"account_id"`
	Limit     int32           `json:"limit"`
	Offset    int32           `json:"offset"`
	From      time.Time       `json:"from"`
	To        time.Time       `json:"to"`
}

func (q *Queries) ListTransactionsByAccountID(ctx context.Context, args ListTxByAccountIDParams) ([]model.Transaction, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> ListTransactionsByAccountID()").Debug()
	rows, err := q.db.QueryContext(ctx, listTXByAccountID, args.AccountID, args.Limit, args.Offset, args.From, args.To)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		q.logs.WithError(err).Warn()
	}()
	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.CategoryID,
			&transaction.Name,
			&transaction.TransactionType,
			&transaction.Amount,
			&transaction.Notes,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.DeletedAt,
		)
		transactions = append(transactions, transaction)
	}
	return transactions, err
}

const listTXByCategoryID = `--name: ListTransactionsByCategoryID :many
SELECT * FROM transactions
WHERE category_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
AND date > $4
AND date < $5
ORDER BY transaction_id
LIMIT  $2
OFFSET $3`

type ListTxByCategoryIDParams struct {
	CategoryID model.CategoryID `json:"category_id"`
	Limit      int32            `json:"limit"`
	Offset     int32            `json:"offset"`
	From       time.Time        `json:"from"`
	To         time.Time        `json:"to"`
}

func (q *Queries) ListTransactionsByCategoryID(ctx context.Context, args ListTxByCategoryIDParams) ([]model.Transaction, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> ListTransactionsByCategoryID()").Debug()
	rows, err := q.db.QueryContext(ctx, listTXByCategoryID, args.CategoryID, args.Limit, args.Offset, args.From, args.To)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		q.logs.WithError(err).Warn()
	}()
	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.CategoryID,
			&transaction.Name,
			&transaction.TransactionType,
			&transaction.Amount,
			&transaction.Notes,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.DeletedAt,
		)
		transactions = append(transactions, transaction)
	}
	return transactions, err
}

const deleteTransaction = `--name: DeleteTransaction :one
UPDATE transactions SET deleted_at = now()
WHERE transaction_id = $1
  AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at`

func (q *Queries) DeleteTransaction(ctx context.Context, id model.TransactionID) (time.Time, error) {
	q.logs.WithField("func", "database/sqlc/transaction.go -> DeleteTransaction()").Debug()
	row := q.db.QueryRowContext(ctx, deleteTransaction, id)
	var transaction model.Transaction
	err := row.Scan(
		&transaction.DeletedAt,
	)
	return transaction.DeletedAt, err
}
