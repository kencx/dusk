package ui

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	p, err := s.db.GetAllBooks(defaultBookFilters())
	if err != nil {
		slog.Error("[ui] failed to load index page", slog.Any("err", err))
		views.NewIndex(s.base, page.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}
	views.NewIndex(s.base, *p, nil).Render(rw, r)
}

func (s *Handler) bookSearch(rw http.ResponseWriter, r *http.Request) {

	filters := initBookFilters(r)
	if errMap := validator.Validate(filters.Filters); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		partials.BookSearchResults(page.Page[dusk.Book]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	p, err := s.db.GetAllBooks(filters)
	if err != nil {
		if err == dusk.ErrNoRows {
			partials.BookSearchResults(page.Page[dusk.Book]{}, err).Render(r.Context(), rw)
			return
		} else {
			slog.Error("failed to get all books", slog.Any("err", err))
			partials.BookSearchResults(page.Page[dusk.Book]{}, err).Render(r.Context(), rw)
			return
		}
	}
	partials.BookSearchResults(*p, nil).Render(r.Context(), rw)
}

func (s *Handler) bookPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(int64(id))
	if err != nil {
		slog.Error("[ui] failed to find book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	if r.URL.Query().Has("delete") {
		response.AddHxTriggerAfterSwap(rw, `{"openModal": ""}`)
		views.DeleteBookModal(book).Render(r.Context(), rw)
		return
	}

	authors, err := s.db.GetAuthorsFromBook(id)
	if err != nil {
		slog.Error("[ui] failed to fetch authors of book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	tags, err := s.db.GetTagsFromBook(id)
	if err != nil {
		slog.Error("[ui] failed to fetch tags of book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	views.NewBook(s.base, book, authors, tags, nil).Render(rw, r)
}

func (s *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteBook(id)
	if err != nil {
		slog.Error("[ui] failed to delete book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}
	// redirect to index page
	response.HxRedirect(rw, r, "/")
}
