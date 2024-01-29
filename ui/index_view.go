package ui

import (
	"dusk"
	"dusk/ui/pages"
	"errors"
	"net/http"
)

func (s *Handler) indexView(rw http.ResponseWriter, r *http.Request) {
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
