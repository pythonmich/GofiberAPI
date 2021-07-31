package database

import (
	model "FiberFinanceAPI/database/models"
	"FiberFinanceAPI/utils"
	"context"
	"time"
)

const createAccount = `INSERT INTO accounts(user_id, account_name, account_type, balance, currency) 
VALUES ($1, $2, $3, $4, $5)
RETURNING account_id, user_id, account_name, account_type, balance, currency, created_at, deleted_at`

type CreateAccountParams struct {
	UserID      model.UserID       `json:"user_id"`
	AccountName string             `json:"account_name"`
	AccountType model.AccountType  `json:"account_type"`
	Balance     int64              `json:"balance"`
	Currency    utils.CurrencyCode `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, args CreateAccountParams) (model.Account, error) {
	q.logs.WithField("func", "database/sqlc/accounts.go -> CreateAccount()").Debug()

	row := q.db.QueryRowContext(ctx, createAccount, args.UserID, args.AccountName, args.AccountType, args.Balance, args.Currency)
	var account model.Account
	err := row.Scan(
		&account.AccountID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.DeletedAt,
	)
	return account, err
}

/*
	we shall only change balance for our user, if our user wants to change type, currency or name we shall
	create a different account in with we will transfer their balance amount if user wants
*/
const updateAccount = `--name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE account_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING account_id, user_id, account_name, account_type, balance, currency, created_at, deleted_at`

type UpdateAccountParams struct {
	AccountID model.AccountID `json:"account_id"`
	Balance   int64           `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, args UpdateAccountParams) (model.Account, error) {
	q.logs.WithField("func", "database/sqlc/accounts.go -> UpdateAccount()").Debug()
	row := q.db.QueryRowContext(ctx, updateAccount, args.AccountID, args.Balance)
	var account model.Account
	err := row.Scan(
		&account.AccountID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.DeletedAt,
	)
	return account, err
}

const getAccountByID = `--name: GetAccountByID :one 
SELECT * FROM accounts
WHERE account_id = $1 
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1`

func (q *Queries) GetAccountByID(ctx context.Context, id model.AccountID) (model.Account, error) {
	q.logs.WithField("func", "database/sqlc/accounts.go -> GetAccountByID()").Debug()

	row := q.db.QueryRowContext(ctx, getAccountByID, id)
	var account model.Account
	err := row.Scan(
		&account.AccountID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.DeletedAt,
	)
	return account, err
}

const listAccounts = `--name: ListAccounts :many
SELECT * FROM accounts
WHERE user_id = $1 
AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY account_id
LIMIT $2
OFFSET $3
`

type ListAccountParams struct {
	UserID model.UserID `json:"user_id"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, args ListAccountParams) ([]model.Account, error) {
	q.logs.WithField("func", "database/sqlc/accounts.go -> ListAccounts()").Debug()
	rows, err := q.db.QueryContext(ctx, listAccounts, args.UserID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			q.logs.WithError(err).Warn("rows not closed")
		}
		q.logs.Debug("rows closed successfully")
	}()
	var accounts []model.Account
	for rows.Next() {
		var account model.Account
		err = rows.Scan(
			&account.AccountID,
			&account.UserID,
			&account.Name,
			&account.Type,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
			&account.DeletedAt,
		)
		accounts = append(accounts, account)
	}
	return accounts, err
}

const deleteAccount = `--name: DeleteAccount :one
UPDATE accounts SET deleted_at = now()
WHERE account_id = $1 
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at`

func (q *Queries) DeleteAccount(ctx context.Context, id model.AccountID) (time.Time, error) {
	q.logs.WithField("func", "database/sqlc/accounts.go -> DeleteAccount()").Debug()
	row := q.db.QueryRowContext(ctx, deleteAccount, id)
	var account model.Account
	err := row.Scan(
		&account.DeletedAt,
	)
	return account.DeletedAt, err
	//return reflect.DeepEqual(account, Account{}) , err
}
