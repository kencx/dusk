package ui

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) authorList(rw http.ResponseWriter, r *http.Request) {
	filters := initSearchFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.AuthorSearchResults(page.Page[dusk.Author]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	authors, err := s.db.GetAllAuthors(filters)
	if err != nil {
		slog.Error("[ui] failed to get all authors", slog.Any("err", err))
		views.NewAuthorList(s.base, page.Page[dusk.Author]{}, filters.Base, err).Render(rw, r)
		return
	}
	views.NewAuthorList(s.base, *authors, filters.Base, nil).Render(rw, r)
}

func (s *Handler) authorSearch(rw http.ResponseWriter, r *http.Request) {
	// If not htmx request, return the full page instead of partial.
	// Required to support hx-push-urls
	if request.IsHtmxRequest(r) {
		s.authorList(rw, r)
		return
	}

	filters := initSearchFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.AuthorSearchResults(page.Page[dusk.Author]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	p, err := s.db.GetAllAuthors(filters)
	if err != nil {
		if err == dusk.ErrNoRows {
			views.AuthorSearchResults(page.Page[dusk.Author]{}, err).Render(r.Context(), rw)
			return
		} else {
			slog.Error("failed to get all authors", slog.Any("err", err))
			views.AuthorSearchResults(page.Page[dusk.Author]{}, err).Render(r.Context(), rw)
			return
		}
	}
	views.AuthorSearchResults(*p, nil).Render(r.Context(), rw)
}

func (s *Handler) authorPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	filters := initBookFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.AuthorSearchResults(page.Page[dusk.Author]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	author, err := s.db.GetAuthor(id)
	if err != nil {
		slog.Error("[ui] failed to get author", slog.Int64("id", id), slog.Any("err", err))
		views.NewAuthor(s.base, dusk.Author{}, page.Page[dusk.Book]{}, filters.Base, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromAuthor(author.Id, filters)
	if err != nil {
		slog.Error("[ui] failed to get books from author", slog.Int64("id", id), slog.Any("err", err))
		views.NewAuthor(s.base, dusk.Author{}, page.Page[dusk.Book]{}, filters.Base, err).Render(rw, r)
		return
	}
	views.NewAuthor(s.base, *author, *books, filters.Base, nil).Render(rw, r)
}
