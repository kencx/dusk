package http

import (
	"dusk"
	"dusk/http/response"
	"dusk/ui/pages"
	"errors"
	"net/http"
)

func (s *Server) rootHandler(rw http.ResponseWriter, r *http.Request) {
	books, err := s.db.GetAllBooks()
	if err != nil && !errors.Is(err, dusk.ErrNoRows) {
		response.InternalServerError(rw, r, err)
		return
	}

	if books == nil {
		books = dusk.Books{}
	}
	pages.Index(books).Render(r.Context(), rw)
}

func (s *Server) bookPageHandler(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(id)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	pages.BookPage(book).Render(r.Context(), rw)
}

func (s *Server) authorPageHandler(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	author, err := s.db.GetAuthor(id)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	pages.AuthorPage(author).Render(r.Context(), rw)
}
