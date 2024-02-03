package ui

import (
	"dusk/http/request"
	"dusk/ui/views"
	"log"
	"net/http"
)

func (s *Handler) bookPage(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(int64(id))
	if err != nil {
		log.Println(err)
		return
	}

	views.BookPage(book).Render(r.Context(), rw)
}
