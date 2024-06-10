package ui

import (
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) tagList(rw http.ResponseWriter, r *http.Request) {
	tags, err := s.db.GetAllTags()
	if err != nil {
		slog.Error("[ui] failed to get all tags", slog.Any("err", err))
		views.NewTagList(s.baseView, nil, err).Render(rw, r)
		return
	}
	views.NewTagList(s.baseView, tags, nil).Render(rw, r)
}

func (s *Handler) tagPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	tag, err := s.db.GetTag(id)
	if err != nil {
		slog.Error("[ui] failed to get tag", slog.Int64("id", id), slog.Any("err", err))
		views.NewTag(s.baseView, nil, nil, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromTag(tag.Id)
	if err != nil {
		slog.Error("[ui] failed to get books from tag", slog.Int64("id", id), slog.Any("err", err))
		views.NewTag(s.baseView, nil, nil, err).Render(rw, r)
		return
	}
	views.NewTag(s.baseView, tag, books, nil).Render(rw, r)
}
