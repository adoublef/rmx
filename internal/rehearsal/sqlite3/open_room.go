package sqlite3

import (
	"context"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

func OpenRoom(ctx context.Context, db sql.DB, r *rehearsal.Room) error {
	var (
		qry = `
		INSERT INTO rooms (id, title, capacity, owner)
		VALUES (?, ?, ?, ?)`
	)

	_, err := db.ExecContext(ctx, qry, r.ID, r.Title, r.Capacity, r.Owner)
	if err != nil {
		return err
	}
	return nil
}
