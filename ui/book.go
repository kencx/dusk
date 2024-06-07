package ui

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	books, err := s.db.GetAllBooks(nil)
	if err != nil {
		slog.Error("[ui] failed to load index page", slog.Any("err", err))
		views.NewIndex(s.baseView, nil, err).Render(rw, r)
		return
	}
	views.NewIndex(s.baseView, books, nil).Render(rw, r)
}

func (s *Handler) bookSearch(rw http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	var input = &dusk.BookFilters{
		Search: readString(qs, "itemSearch", ""),
		Title:  readString(qs, "title", ""),
		Author: readString(qs, "author", ""),
	}

	books, err := s.db.GetAllBooks(input)
	if err != nil {
		log.Println(err)
		partials.BookSearchResults(nil, err).Render(r.Context(), rw)
		return
	}
	partials.BookSearchResults(books, nil).Render(r.Context(), rw)
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
		views.NewBook(s.baseView, nil, err).Render(rw, r)
		return
	}

	if r.URL.Query().Has("delete") {
		response.AddHxTriggerAfterSwap(rw, `{"openModal": ""}`)
		views.DeleteBookModal(book).Render(r.Context(), rw)
		return
	}
	views.NewBook(s.baseView, book, nil).Render(rw, r)
}

func (s *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteBook(id)
	if err != nil {
		slog.Error("[ui] failed to delete book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.baseView, nil, err).Render(rw, r)
		return
	}
	// redirect to index page
	response.HxRedirect(rw, r, "/")
}
