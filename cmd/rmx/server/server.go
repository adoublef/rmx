package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	service "github.com/rapidmidiex/rmx/internal/http"
	sql "github.com/rapidmidiex/rmx/sqlite3"
)

type Server struct {
	s *http.Server
}

// ListenAndServe listens on the TCP network address and then handles requests on incoming connections.
func (s *Server) ListenAndServe() error {
	err := s.s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http listen and serve: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.s.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http shutdown: %w", err)
	}
	return nil
}

// A New Server defines parameters for running an HTTP server.
func New(addr string, db sql.DB) (*Server, error) {
	mux := service.New(db)
	s := &http.Server{Addr: addr, Handler: mux}
	return &Server{s}, nil
}
