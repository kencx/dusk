package ui

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

// Render index page and book library
func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	filters := initBookFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.NewIndex(s.base, page.Page[dusk.Book]{}, filters.Base, errors.New("validation error")).Render(rw, r)
		return
	}

	p, err := s.db.GetAllBooks(filters)
	if err != nil {
		slog.Error("[ui] failed to load index page", slog.Any("err", err))
		views.NewIndex(s.base, page.Page[dusk.Book]{}, filters.Base, err).Render(rw, r)
		return
	}

	views.NewIndex(s.base, *p, filters.Base, nil).Render(rw, r)
}

func (s *Handler) notFound(rw http.ResponseWriter, r *http.Request) {
	s.base.NotFound().Render(r.Context(), rw)
}
