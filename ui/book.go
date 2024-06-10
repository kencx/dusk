package ui

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	page, err := s.db.GetAllBooks(dusk.DefaultBookFilters())
	if err != nil {
		slog.Error("[ui] failed to load index page", slog.Any("err", err))
		views.NewIndex(s.base, dusk.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}
	views.NewIndex(s.base, *page, nil).Render(rw, r)
}

func (s *Handler) bookSearch(rw http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	var filters = &dusk.BookFilters{
		Title:  readString(qs, "title", ""),
		Author: readString(qs, "author", ""),
		SearchFilters: dusk.SearchFilters{
			Search: readString(qs, "itemSearch", ""),
			Filters: dusk.Filters{
				AfterId:      readInt(qs, "after_id", 0),
				PageSize:     readInt(qs, "page_size", 30),
				Sort:         readString(qs, "sort", "title"),
				SortSafeList: dusk.DefaultSafeList(),
			},
		},
	}

	if errMap := validator.Validate(filters.Filters); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		partials.BookSearchResults(dusk.Page[dusk.Book]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	page, err := s.db.GetAllBooks(filters)
	if err != nil {
		if err == dusk.ErrNoRows {
			partials.BookSearchResults(dusk.Page[dusk.Book]{}, err).Render(r.Context(), rw)
			return
		} else {
			slog.Error("failed to get all books", slog.Any("err", err))
			partials.BookSearchResults(dusk.Page[dusk.Book]{}, err).Render(r.Context(), rw)
			return
		}
	}

	// only return page partial when not querying for first page
	if !page.First() {
		partials.BookPage(*page).Render(r.Context(), rw)
	} else {
		partials.BookSearchResults(*page, nil).Render(r.Context(), rw)
	}
}

// read int query parameter
func readInt(qv url.Values, key string, defaultValue int) int {
	value := qv.Get(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

// read string query parameter
func readString(qv url.Values, key string, defaultValue string) string {
	value := qv.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (s *Handler) bookPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(int64(id))
	if err != nil {
		slog.Error("[ui] failed to find book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, err).Render(rw, r)
		return
	}

	if r.URL.Query().Has("delete") {
		response.AddHxTriggerAfterSwap(rw, `{"openModal": ""}`)
		views.DeleteBookModal(book).Render(r.Context(), rw)
		return
	}
	views.NewBook(s.base, book, nil).Render(rw, r)
}

func (s *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteBook(id)
	if err != nil {
		slog.Error("[ui] failed to delete book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, err).Render(rw, r)
		return
	}
	// redirect to index page
	response.HxRedirect(rw, r, "/")
}
