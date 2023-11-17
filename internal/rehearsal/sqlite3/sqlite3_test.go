package sqlite3_test

import (
	"context"
	"embed"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
	testsql "github.com/rapidmidiex/rmx/sqlite3"
	"github.com/rs/xid"
)

func TestSqlite3(t *testing.T) {
	t.Run("open a single room", withClient(func(is *is.I, db testsql.DB) {
		var (
			ctx  = context.Background()
			room = rehearsal.NewRoom("Backtrack Boys", 4, xid.New())
		)

		err := sql.OpenRoom(ctx, db, room)
		is.NoErr(err) // open a new room using valid credentials
	}))

	t.Run("opening a room with invalid capacity", withClient(func(is *is.I, db testsql.DB) {
		var (
			ctx  = context.Background()
			room = rehearsal.Room{ID: xid.New(), Title: "The Foo Fighters", Capacity: 0, Owner: xid.New()}
		)

		err := sql.OpenRoom(ctx, db, &room)
		is.True(err != nil) // invalid room capacity
	}))

	// TODO test that invalid characters in room name throws an error
}

//go:embed all:*.up.sql
var up embed.FS

func withClient(f func(is *is.I, db testsql.DB)) func(t *testing.T) {
	var ctx = context.Background()

	return func(t *testing.T) {
		// run in parallel

		dsn := filepath.Join(t.TempDir(), "test.db")
		db, err := testsql.Open(dsn)
		if err != nil {
			t.Fatalf("open connection: %v", err)
		}
		defer db.Close()
		if err := testsql.Up(ctx, db, up); err != nil {
			t.Fatalf("run migrations: %v", err)
		}
		f(is.New(t), db)
	}
}
