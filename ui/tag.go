package ui

import (
	"log"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) tagList(rw http.ResponseWriter, r *http.Request) {
	tags, err := s.db.GetAllTags()
	if err != nil {
		log.Println(err)
		views.NewTagList(nil, err).Render(rw, r)
		return
	}
	views.NewTagList(tags, nil).Render(rw, r)
}

func (s *Handler) tagPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	tag, err := s.db.GetTag(id)
	if err != nil {
		log.Println(err)
		views.NewTag(nil, nil, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromTag(tag.Id)
	if err != nil {
		log.Println(err)
		views.NewTag(nil, nil, err).Render(rw, r)
		return
	}
	views.NewTag(tag, books, nil).Render(rw, r)
}
