package ui

import (
	"log"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) authorList(rw http.ResponseWriter, r *http.Request) {
	authors, err := s.db.GetAllAuthors()
	if err != nil {
		log.Println(err)
		views.NewAuthorList(nil, err).Render(rw, r)
		return
	}
	views.NewAuthorList(authors, nil).Render(rw, r)
}

func (s *Handler) authorPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	author, err := s.db.GetAuthor(id)
	if err != nil {
		log.Println(err)
		views.NewAuthor(nil, nil, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromAuthor(author.ID)
	if err != nil {
		log.Println(err)
		views.NewAuthor(nil, nil, err).Render(rw, r)
		return
	}

	// handle toggle
	if r.URL.Query().Has("show") {
		show := partials.LibraryView(r.URL.Query().Get("show"))
		partials.ViewToggle(books, show).Render(r.Context(), rw)
		return
	}

	views.NewAuthor(author, books, nil).Render(rw, r)
}
