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

func (s *Server) indexView(rw http.ResponseWriter, r *http.Request) {
	var errMsg string

	books, err := s.db.GetAllBooks()
	if err != nil {
		switch {
		case errors.Is(err, dusk.ErrNoRows):
			errMsg = dusk.ErrNoRows.Error()
		default:
			errMsg = "Something went wrong."
		}
	}

	if books == nil {
		books = dusk.Books{}
	}

	pages.Index(books, errMsg).Render(r.Context(), rw)
}

func (s *Server) importView(rw http.ResponseWriter, r *http.Request) {
	// handle tabs
	if r.URL.Query().Has("tab") {
		tab := r.URL.Query().Get("tab")
		pages.Tab(tab).Render(r.Context(), rw)
		return
	}

	pages.Import("openlibrary", "").Render(r.Context(), rw)
}

func (s *Server) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("openlibrary")
	m, err := metadata.Fetch(isbn)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		pages.Import("openlibrary", "Something went wrong").Render(r.Context(), rw)
		return
	}

	b := m.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		s.ErrorLog.Printf("err: %v", err)
		pages.Import("openlibrary", "Something went wrong").Render(r.Context(), rw)
		return
	}

	result, err := s.db.CreateBook(b)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		pages.Import("openlibrary", "Something went wrong").Render(r.Context(), rw)
		return
	}

	s.InfoLog.Printf("New book created: %v", result)
	pages.Import("openlibrary", "Book imported").Render(r.Context(), rw)
}

func (s *Server) bookView(rw http.ResponseWriter, r *http.Request) {
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

func (s *Server) authorView(rw http.ResponseWriter, r *http.Request) {
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
