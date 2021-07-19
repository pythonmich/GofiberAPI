package database

import (
	"FiberFinanceAPI/utils"
	"database/sql"
)

type Repo interface {
	QueryInterface
}

type SQLRepo struct {
	*Queries
	db *sql.DB
}

func NewRepo(db *sql.DB, log *utils.StandardLogger) Repo {
	return SQLRepo{
		Queries: NewQuery(db, log),
		db:      db,
	}
}
