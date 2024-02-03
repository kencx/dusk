package ui

import (
	"dusk"
	"dusk/ui/views"
	"errors"
	"log"
	"net/http"
)

func (s *Handler) authorList(rw http.ResponseWriter, r *http.Request) {
	// m := views.NewAuthorListViewModel(nil, nil)

	authors, err := s.db.GetAllAuthors()
	if err != nil {
		switch {
		case errors.Is(err, dusk.ErrNoRows):
			log.Println(err)
			// TODO set custom message
			// m.RenderError(rw, r, err)
		default:
			log.Println(err)
			// m.RenderError(rw, r, err)
		}
		return
	}

	if authors == nil {
		authors = dusk.Authors{}
	}

	views.AuthorList(authors).Render(r.Context(), rw)
	// m.Authors = authors
	// m.Render(rw, r)
}
