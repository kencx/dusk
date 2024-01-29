package ui

import (
	"dusk"
	"dusk/ui/pages"
	"errors"
	"net/http"
)

func (s *Handler) indexView(rw http.ResponseWriter, r *http.Request) {
	m := pages.NewIndexViewModel(nil, nil)

	books, err := s.db.GetAllBooks()
	if err != nil {
		switch {
		case errors.Is(err, dusk.ErrNoRows):
			// TODO set custom message
			m.RenderError(rw, r, err)
		default:
			m.RenderError(rw, r, err)
		}
		return
	}

	if books == nil {
		books = dusk.Books{}
	}
	m.Books = books
	m.Render(rw, r)
}
