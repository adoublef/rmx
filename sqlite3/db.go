package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/maragudk/migrate"
	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

var args = strings.Join([]string{"_journal=wal", "_timeout=5000", "_synchronous=normal", "_fk=true"}, "&")

type DB interface {
	Close() error
	Exec(ctx context.Context, query string, args ...any) (rowsAffected int64, err error)
	Query(ctx context.Context, query string, args ...any) (ScanIterator, error)
	QueryRow(ctx context.Context, query string, args ...any) Scanner
}

var _ DB = (*Conn)(nil)

type Conn struct {
	rwc *sql.DB
}

// Query implements DB.
func (c *Conn) Query(ctx context.Context, query string, args ...any) (ScanIterator, error) {
	rs, err := c.rwc.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return rs, nil
}

// Exec executes a query without returning any rows.
func (c *Conn) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	rs, err := c.rwc.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("execute: %w", err)
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	return n, nil
}

// QueryRow executes a query that is expected to return at most one row.
func (c *Conn) QueryRow(ctx context.Context, query string, args ...any) Scanner {
	return c.rwc.QueryRowContext(ctx, query, args...)
}

// Close closes the database and prevents new queries from starting. C
func (c *Conn) Close() error {
	return c.rwc.Close()
}

// Open opens a database connection for the given sqlite file.
func Open(dsn string) (*Conn, error) {
	db, err := sql.Open(driverName, dsn+"?"+args)
	if err != nil {
		return nil, fmt.Errorf("open sqlite3: %w", err)
	}
	return &Conn{db}, nil
}

type Scanner interface {
	Err() error
	Scan(dest ...any) error
}

type ScanIterator interface {
	io.Closer
	Scanner
	Next() bool
}

// Up from the current version.
func Up(ctx context.Context, db *Conn, fsys fs.FS) error {
	err := migrate.Up(ctx, db.rwc, fsys)
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
