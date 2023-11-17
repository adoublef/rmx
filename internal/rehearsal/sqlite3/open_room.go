package sqlite3

import (
	"context"
	"embed"
	"fmt"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

//go:embed all:*.up.sql
var FS embed.FS

func OpenRoom(ctx context.Context, db sql.DB, r rehearsal.Room) error {
	var (
		qry = `
		INSERT INTO rooms (id, title, capacity, owner)
		VALUES (?, ?, ?, ?)`
	)

	_, err := db.Exec(ctx, qry, r.ID, r.Title, r.Capacity, r.Owner)
	if err != nil {
		return fmt.Errorf("open room: %w", err)
	}
	return nil
}
