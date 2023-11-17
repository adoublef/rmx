package sqlite3

import (
	"context"
	"fmt"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

func ListRooms(ctx context.Context, db sql.DB, limit, offset int) ([]rehearsal.Room, int, error) {
	var (
		qry = `
		SELECT r.id, r.title, r.capacity, r.owner
		FROM rooms r
		LIMIT ? OFFSET ?`
	)

	rs, err := db.Query(ctx, qry, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list rooms: %w", err)
	}

	var vv []rehearsal.Room
	for rs.Next() {
		var v rehearsal.Room
		err := rs.Scan(&v.ID, &v.Title, &v.Capacity, &v.Owner)
		if err != nil {
			return nil, 0, fmt.Errorf("list rooms: %w", err)
		}
		vv = append(vv, v)
	}
	if err := rs.Err(); err != nil {
		return nil, 0, fmt.Errorf("list rooms: %w", err)
	}

	return vv, len(vv), nil
}
