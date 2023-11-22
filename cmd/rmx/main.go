package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rapidmidiex/rmx/cmd/rmx/server"
	"github.com/rapidmidiex/rmx/errgroup"
	room "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

const (
	E = flag.ExitOnError
	S = "s"
	M = "m"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("too few arguments")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// serve
	s := serve{fs: flag.NewFlagSet(S, E)}
	s.fs.StringVar(&s.addr, "addr", ":8080", "listen address")
	s.fs.StringVar(&s.dsn, "dsn", ":memory:", "datasource name")
	// migrate
	m := migrate{fs: flag.NewFlagSet(M, E)}
	m.fs.StringVar(&m.dsn, "dsn", ":memory:", "datasource name")

	var err error
	switch sub, args := os.Args[1], os.Args[2:]; sub {
	case S:
		err = s.serve(ctx, args)
	case M:
		err = m.migrate(ctx, args)
	default:
		err = errors.New("unknown subcommand")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// serve
type serve struct {
	fs   *flag.FlagSet
	addr string
	dsn  string
}

func (c *serve) serve(ctx context.Context, args []string) (err error) {
	if err = c.fs.Parse(args); err != nil {
		return
	}

	var (
		eg = errgroup.New(ctx)
	)

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
	fs  *flag.FlagSet
	dsn string
}

func (c *migrate) migrate(ctx context.Context, args []string) (err error) {
	if err = c.fs.Parse(args); err != nil {
		return
	}

	var (
		eg = errgroup.New(ctx)
	)

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	eg.Go(func(ctx context.Context) error {
		return room.Up(ctx, db)
	})

	return eg.Wait()
}
