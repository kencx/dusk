package ui

import (
	"dusk"
	"dusk/ui/views"
	"errors"
	"net/http"
)

func (s *Handler) index(rw http.ResponseWriter, r *http.Request) {
	m := views.NewIndexViewModel(nil, nil)

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
