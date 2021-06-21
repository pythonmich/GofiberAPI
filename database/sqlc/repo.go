package database

import "database/sql"

type Repo interface {
	QueryInterface
}

type SQLRepo struct {
	*Queries
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return SQLRepo{
		Queries: NewQuery(db),
		db:      db,
	}
}


