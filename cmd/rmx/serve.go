package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/choria-io/fisk"
	"github.com/rapidmidiex/rmx/cmd/rmx/server"
	eg "github.com/rapidmidiex/rmx/errgroup"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

const (
	serveName = "serve"
	serveHelp = "Run application server"
	addrName  = "addr"
	addrHelp  = "Listen address"
	dsnName   = "dsn"
	dsnShort  = "d"
	dsnHelp   = "Datasource name"
)

type serveCmd struct {
	addr string
	dsn  string
}

func configServeCmd(a app) {
	c := &serveCmd{}
	serve := a.Command(serveName, serveHelp).Action(c.serve)
	serve.Flag(addrName, addrHelp).StringVar(&c.addr)
	serve.Flag(dsnName, dsnHelp).StringVar(&c.dsn)
}

func init() { setCmd("serve", 0, configServeCmd) }

func (c *serveCmd) serve(_ *fisk.ParseContext) error {
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
	
	_, eg, cancel := newErrgroup()
	defer cancel()

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

func newErrgroup() (context.Context, *eg.Group, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	return ctx, eg.New(ctx), cancel
}
