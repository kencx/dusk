package ui

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) authorList(rw http.ResponseWriter, r *http.Request) {
	authors, err := s.db.GetAllAuthors(nil)
	if err != nil {
		slog.Error("[ui] failed to get all authors", slog.Any("err", err))
		views.NewAuthorList(s.baseView, nil, err).Render(rw, r)
		return
	}
	views.NewAuthorList(s.baseView, authors, nil).Render(rw, r)
}

func (s *Handler) authorSearch(rw http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	input := &dusk.SearchFilters{
		Search: readString(qs, "itemSearch", ""),
	}

	authors, err := s.db.GetAllAuthors(input)
	if err != nil {
		log.Println(err)
		views.AuthorSearchResults(nil, err).Render(r.Context(), rw)
		return
	}
	views.AuthorSearchResults(authors, nil).Render(r.Context(), rw)
}

func (s *Handler) authorPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	author, err := s.db.GetAuthor(id)
	if err != nil {
		slog.Error("[ui] failed to get author", slog.Int64("id", id), slog.Any("err", err))
		views.NewAuthor(s.baseView, nil, nil, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromAuthor(author.Id)
	if err != nil {
		slog.Error("[ui] failed to get books from author", slog.Int64("id", id), slog.Any("err", err))
		views.NewAuthor(s.baseView, nil, nil, err).Render(rw, r)
		return
	}
	views.NewAuthor(s.baseView, author, &dusk.BooksPage{
		Page: dusk.Page{
			Total: int64(len(books)),
		},
		Books: books,
	}, nil).Render(rw, r)
}
