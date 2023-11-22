package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rapidmidiex/rmx/internal/rehearsal"
	sql "github.com/rapidmidiex/rmx/internal/rehearsal/sqlite3"
	"github.com/rs/xid"
)

func (s *Service) handleOpenRoom() http.HandlerFunc {
	var parseRoom = func(r *http.Request, uid xid.ID) (rehearsal.Room, error) {
		var (
			title = r.PostFormValue("title")
			cap   = r.PostFormValue("capacity")
			// todo - xsrf
		)

		c, err := strconv.Atoi(cap)
		if err != nil {
			return rehearsal.Room{}, fmt.Errorf("string convert: %w", err)
		}

		return rehearsal.ParseRoom(title, c, uid)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx      = r.Context()
			redirect = "/"
			// note - for now will use session id
			uid = xid.New()
		)

		rm, err := parseRoom(r, uid)
		if err != nil {
			http.Error(w, "Invalid form details", http.StatusUnprocessableEntity)
			return
		}

		if err = sql.OpenRoom(ctx, s.db, rm); err != nil {
			http.Error(w, "Failed to open room", http.StatusInternalServerError)
			return
		}

		// 303
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}
