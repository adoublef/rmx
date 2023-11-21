package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/choria-io/fisk"
	"github.com/rapidmidiex/rmx/cmd/rmx/server"
	"github.com/rapidmidiex/rmx/errgroup"
	room "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

const (
	appName     = "audio"
	appHelp     = ""
	serveName   = "serve"
	serveHelp   = "Run application"
	migrateName = "migrate"
	migrateHelp = "Run database migrations"
	addrName    = "addr"
	addrHelp    = "HTTP Listen address"
	dsnName     = "dsn"
	dsnShort    = 'd'
	dsnHelp     = "Datasource name"
)

func main() {
	f := fisk.New(appName, appHelp)
	// serve
	{
		v := &serve{}
		s := f.Command(serveName, serveHelp).Action(v.serve)
		s.Flag(addrName, addrHelp).StringVar(&v.addr)
		s.Flag(dsnName, dsnHelp).Short(dsnShort).StringVar(&v.dsn)
	}
	// migrate
	{
		v := &migrate{}
		m := f.Command(migrateName, migrateHelp).Action(v.migrate)
		m.Flag(dsnName, dsnHelp).Short(dsnShort).StringVar(&v.dsn)
	}
	// parse flags
	_, err := f.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

// serve
type serve struct {
	addr string
	dsn  string
}

func (c *serve) serve(_ *fisk.ParseContext) error {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		eg          = errgroup.New(ctx)
	)
	defer cancel()

	if c.addr == "" {
		c.addr = ":8080"
	}
	if c.dsn == "" {
		c.dsn = ":memory:"
	}

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	
	s, err := server.New(c.addr, db)
	if err != nil {
		return err
	}

	eg.Go(func(ctx context.Context) error {
		return s.ListenAndServe()
	})
	eg.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return s.Shutdown()
	})

	return eg.Wait()
}

// migrate
type migrate struct {
	dsn string
}

func (c *migrate) migrate(_ *fisk.ParseContext) error {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	)
	defer cancel()

	if c.dsn == "" {
		c.dsn = ":memory:"
	}

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return sql.Up(ctx, db, room.FS)
}
