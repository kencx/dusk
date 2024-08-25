package ui

import (
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/ui/views"
)

// Render index page and book library
func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	p, err := s.db.GetAllBooks(defaultBookFilters())
	if err != nil {
		slog.Error("[ui] failed to load index page", slog.Any("err", err))
		views.NewIndex(s.base, page.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}
	views.NewIndex(s.base, *p, nil).Render(rw, r)
}

func (s *Handler) notFound(rw http.ResponseWriter, r *http.Request) {
	s.base.NotFound().Render(r.Context(), rw)
}
