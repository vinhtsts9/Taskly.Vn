package database

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error // Thêm dòng này
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
