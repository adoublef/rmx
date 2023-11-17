package http

import (
	"net/http"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
)

func (s *Service) handleIndex() http.HandlerFunc {
	type response struct {
		Rooms []*rehearsal.Room
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
		)
		// before rendering the page, check available rooms
		// include limits and cursor based queries
		rr, _, err := sql.ListRooms(ctx, s.db, 10, 0)
		if err != nil {
			http.Error(w, "Cannot list rooms", http.StatusInternalServerError)
			return
		}

		v := &response{Rooms: rr}

		t, _ := s.fs.ParseFiles(pageIndex)
		t.Execute(w, v)
	}
}
