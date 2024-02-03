package ui

import (
	"dusk/http/request"
	"log"
	"net/http"
)

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
