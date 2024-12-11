package prizzle

import (
	"context"
	"database/sql"
)

type DatabaseClient interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	NewQuery() *SqlQuery
}

type Transactor interface {
	DatabaseClient
	Commit() error
	Rollback() error
}
