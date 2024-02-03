package ui

import (
	"dusk/http/request"
	"dusk/ui/views"
	"log"
	"net/http"
)

func (s *Handler) authorPage(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	author, err := s.db.GetAuthor(id)
	if err != nil {
		log.Println(err)
		return
	}

	books, err := s.db.GetAllBooksFromAuthor(author.ID)
	if err != nil {
		log.Println(err)
		return
	}

	views.AuthorPage(author, books, "").Render(r.Context(), rw)
}
