package database

import (
	"FiberFinanceAPI/utils"
	"context"
	"database/sql"
)

type DbTx interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type Queries struct {
	db   DbTx
	logs *utils.StandardLogger
}

func NewQuery(db DbTx, log *utils.StandardLogger) *Queries {
	return &Queries{
		db:   db,
		logs: log,
	}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
