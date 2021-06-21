package database

import (
	"context"
	"database/sql"
)

type DbTx interface {
	ExecContext(ctx context.Context, query string, args ...interface{})(sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{})*sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{})(*sql.Rows, error)
	PrepareContext(ctx context.Context, query string)(*sql.Stmt, error)
}

type Queries struct {
	db DbTx
}


func NewQuery(db DbTx) *Queries {
	return &Queries{db: db}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
