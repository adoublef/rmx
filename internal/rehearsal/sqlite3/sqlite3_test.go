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
			room = rehearsal.NewRoom("Bit Rockers", 4, xid.New())
		)

		err := sql.OpenRoom(ctx, db, room)
		is.NoErr(err) // open a new room using valid credentials
	}))

	t.Run("opening a room with invalid capacity", withClient(func(is *is.I, db testsql.DB) {
		var (
			ctx  = context.Background()
			room = rehearsal.Room{ID: xid.New(), Title: "The Bug Fixes", Capacity: 0, Owner: xid.New()}
		)

		err := sql.OpenRoom(ctx, db, &room)
		is.True(err != nil) // invalid room capacity
	}))

	// TODO test that invalid characters in room name throws an error

	t.Run("open a room and find by id", withClient(func(is *is.I, db testsql.DB) {
		var (
			ctx  = context.Background()
			room = rehearsal.NewRoom("Binary Solo", 1, xid.New())
		)

		err := sql.OpenRoom(ctx, db, room)
		is.NoErr(err) // open new room

		found, err := sql.FindRoom(ctx, db, room.ID)
		is.NoErr(err) // find room using xid
		is.Equal(found, room)
	}))

	t.Run("open a room and fail to find by id", withClient(func(is *is.I, db testsql.DB) {
		var (
			ctx  = context.Background()
			room = rehearsal.NewRoom("Binary Solo", 1, xid.New())
		)

		err := sql.OpenRoom(ctx, db, room)
		is.NoErr(err) // open new room

		_, err = sql.FindRoom(ctx, db, xid.New())
		is.True(err != nil) // failed to find room
	}))

	t.Run("open rooms and list all rooms", withClient(func(is *is.I, db testsql.DB) {
		var (
			ctx = context.Background()
			r1  = rehearsal.NewRoom("8-Bit Maestros", 8, xid.New())
			r2  = rehearsal.NewRoom("Ruby on Rails & Rhythms", 3, xid.New())
			r3  = rehearsal.NewRoom("The Firewall Quartet", 4, xid.New())
		)

		err := sql.OpenRoom(ctx, db, r1)
		is.NoErr(err) // open room 1

		err = sql.OpenRoom(ctx, db, r2)
		is.NoErr(err) // open room 2

		err = sql.OpenRoom(ctx, db, r3)
		is.NoErr(err) // open room 3

		_, n, err := sql.ListRooms(ctx, db, 10, 1)
		is.NoErr(err) // list rooms
		is.Equal(n, 2) // len(rooms) == 2
	}))
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
