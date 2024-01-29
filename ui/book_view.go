package ui

import (
	"dusk/http/request"
	"dusk/ui/pages"
	"log"
	"net/http"
)

func (s *Handler) bookView(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(int64(id))
	if err != nil {
		log.Println(err)
		return
	}

	pages.BookPage(book).Render(r.Context(), rw)
}
