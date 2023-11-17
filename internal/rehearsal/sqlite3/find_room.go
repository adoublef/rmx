package sqlite3

import (
	"context"
	"fmt"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/sqlite3"
	"github.com/rs/xid"
)

func FindRoom(ctx context.Context, db sql.DB, id xid.ID) (*rehearsal.Room, error) {
	var (
		qry = `
		SELECT r.id, r.title, r.capacity, r.owner
		FROM rooms r
		WHERE id = ?`
		r rehearsal.Room
	)

	err := db.QueryRowContext(ctx, qry, id).Scan(
		&r.ID, &r.Title, &r.Capacity, &r.Owner)
	if err != nil {
		return nil,  fmt.Errorf("find room: %w", err)
	}
	return &r, nil
}
