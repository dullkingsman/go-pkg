package prizzle

import (
	"context"
	"database/sql"
)

type DatabaseClientBase interface {
	NewQuery() *SqlQuery
}

type DatabaseClientContextOnly interface {
	DatabaseClientBase
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DatabaseClientContextFree interface {
	DatabaseClientBase
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type DatabaseClient interface {
	DatabaseClientContextOnly
	DatabaseClientContextFree
}

type Transactor interface {
	DatabaseClient
	Commit() error
	Rollback() error
	Stmt(query string) *sql.Stmt
	StmtContext(ctx context.Context, query string) *sql.Stmt
}
