package sqlite3

import (
	"context"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)


func ListRooms(ctx context.Context, db sql.DB) ([]*rehearsal.Room, error) {
	return nil, nil
}