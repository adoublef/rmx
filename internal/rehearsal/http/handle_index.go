package http

import (
	"net/http"
	"strconv"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
)

func (s *Service) handleIndex() http.HandlerFunc {
	type response struct {
		Rooms []*rehearsal.Room
	}
	var parseQueries = func(r *http.Request) (limit, offset int) {
		var (
			// note - should not use these names directly
			// maybe use cursor instead
			l   = r.URL.Query().Get("limit")
			off = r.URL.Query().Get("offset")
			err error
		)

		limit, err = strconv.Atoi(l)
		if err != nil {
			return 10, 0
		}
		offset, err = strconv.Atoi(off)
		if err != nil {
			return 10, 0
		}

		return limit, offset
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = r.Context()
			l, off = parseQueries(r)
		)
		// check url queries for values
		rr, _, err := sql.ListRooms(ctx, s.db, l, off)
		if err != nil {
			http.Error(w, "Cannot list rooms", http.StatusInternalServerError)
			return
		}

		v := &response{Rooms: rr}

		t, _ := s.fs.ParseFiles(pageIndex)
		t.Execute(w, v)
	}
}
