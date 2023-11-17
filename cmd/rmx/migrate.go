package main

import (
	"context"

	"github.com/choria-io/fisk"

	room "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

const (
	migrateName = "migrate"
	migrateHelp = "Run database migrations"
)

type migrateCmd struct {
	dsn string
}

func configMigrateCmd(a app) {
	c := &migrateCmd{}
	serve := a.Command(migrateName, migrateHelp).Action(c.migrate)
	serve.Flag(dsnName, dsnHelp).StringVar(&c.dsn)
}

func init() { setCmd("migrate", 1, configMigrateCmd) }

func (c *migrateCmd) migrate(_ *fisk.ParseContext) error {
	if c.dsn == "" {
		c.dsn = ":memory:"
	}

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return sql.Up(context.Background(), db, room.FS)
}
