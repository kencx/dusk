package ui

import (
	"dusk/http/request"
	"dusk/ui/views"
	"log"
	"net/http"
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
	id := request.HandleInt64("id", rw, r)
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
	views.NewAuthor(author, books, nil).Render(rw, r)
}
