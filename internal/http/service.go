package http

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	sql "github.com/rapidmidiex/rmx/sqlite3"
	"github.com/rapidmidiex/rmx/template"
)

type Service struct {
	m  *chi.Mux
	fs *template.FS
	db sql.DB
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

// A New Service is constructed
func New(db sql.DB) *Service {
	s := Service{chi.NewMux(), FS, db}
	s.routes()
	return &s
}

func (s Service) routes() {
	s.m.Get("/", s.handleIndex())
	s.m.Post("/rooms", s.handleOpenRoom())
	s.m.Get("/rooms/{id}", s.handleEnterRoom())
}

var (
	//go:embed all:*.html
	embedFS   embed.FS
	FS        = template.NewFS(embedFS)
	pageIndex = "index.html"
	pageRoom  = "room.html"
)
