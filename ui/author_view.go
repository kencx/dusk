package ui

import (
	"dusk/http/request"
	"dusk/http/response"
	"dusk/ui/pages"
	"net/http"
)

func (s *Handler) authorView(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
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
