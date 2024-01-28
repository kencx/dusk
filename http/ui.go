package http

import (
	"dusk"
	"dusk/http/response"
	"dusk/metadata"
	"dusk/ui/pages"
	"dusk/validator"
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

func (s *Server) importPageHandler(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		pages.Import().Render(r.Context(), rw)
		return
	}

	if r.Method == http.MethodPost {
		isbn := r.FormValue("openlibrary")
		m, err := metadata.Fetch(isbn)
		if err != nil {
			s.ErrorLog.Printf("err: %v", err)
			response.InternalServerError(rw, r, err)
			return
		}

		b := m.ToBook()

		v := validator.New()
		b.Validate(v)
		if !v.Valid() {
			response.ValidationError(rw, r, v.Errors)
			return
		}

		result, err := s.db.CreateBook(b)
		if err != nil {
			s.ErrorLog.Printf("err: %v", err)
			response.BadRequest(rw, r, err)
			return
		}

		s.InfoLog.Printf("New book created: %v", result)
		http.Redirect(rw, r, "/", http.StatusFound)
	}
}
