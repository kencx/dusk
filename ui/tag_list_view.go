package ui

import (
	"dusk"
	"errors"
	"net/http"
)

func (s *Handler) tagList(rw http.ResponseWriter, r *http.Request) {
	// m := views.NewTagListViewModel(nil, nil)

	tags, err := s.db.GetAllTags()
	if err != nil {
		switch {
		case errors.Is(err, dusk.ErrNoRows):
			// TODO set custom message
			// m.RenderError(rw, r, err)
		default:
			// m.RenderError(rw, r, err)
		}
		return
	}

	if tags == nil {
		tags = dusk.Tags{}
	}
	// m.Tags = tags
	// m.Render(rw, r)
}
