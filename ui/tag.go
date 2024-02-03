package ui

import (
	"dusk"
	"dusk/http/request"
	"errors"
	"log"
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

func (s *Handler) tagPage(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	_, err := s.db.GetTag(int64(id))
	if err != nil {
		log.Println(err)
		return
	}

	// views.TagPage(tag, books, "").Render(r.Context(), rw)
}
