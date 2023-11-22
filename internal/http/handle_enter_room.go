package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
	"github.com/rs/xid"
)

func (s *Service) handleEnterRoom() http.HandlerFunc {
	type response struct {
		Room  rehearsal.Room
	}
	var parseParams = func(r *http.Request) (xid.ID, error) {
		return xid.FromString(chi.URLParam(r, "id"))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx     = r.Context()
			id, err = parseParams(r)
		)
		if err != nil {
			http.Error(w, "Bad id", http.StatusBadRequest)
			return
		}

		found, err := sql.FindRoom(ctx, s.db, id)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// render page
		t, _ := s.fs.ParseFiles(pageIndex, pageRoom)
		t.Execute(w, response{Room: found})
	}
}
