package sqlite3

import (
	"context"
	"embed"

	sql "github.com/rapidmidiex/rmx/sqlite3"
)

//go:embed all:*.up.sql
var fsys embed.FS

var Up = func(ctx context.Context, db *sql.Conn) error {
	return sql.Up(ctx, db, fsys)
}