package database

import (
	"database/sql"
)

type Store interface {
	Querier
}

type SQLStore struct {
	db       *sql.DB
	*Queries // embed Queries để dùng luôn tất cả query bình thường
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
