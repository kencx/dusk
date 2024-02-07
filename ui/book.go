package ui

import (
	"dusk/http/request"
	"dusk/http/response"
	"dusk/ui/views"
	"log"
	"net/http"
)

func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	books, err := s.db.GetAllBooks()
	if err != nil {
		log.Println(err)
		views.NewIndex(nil, err).Render(rw, r)
		return
	}
	views.NewIndex(books, nil).Render(rw, r)
}

func (s *Handler) bookPage(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(int64(id))
	if err != nil {
		log.Println(err)
		views.NewBook(nil, err).Render(rw, r)
		return
	}
	views.NewBook(book, nil).Render(rw, r)
}

func (s *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteBook(id)
	if err != nil {
		log.Println(err)
		views.NewBook(nil, err).Render(rw, r)
		return
	}

	log.Printf("book %d deleted", id)

	// redirect to index page
	// TODO show deleted message
	response.HxRedirect(rw, r, "/")
}
