package ui

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) tagList(rw http.ResponseWriter, r *http.Request) {
	tags, err := s.db.GetAllTags(nil)
	if err != nil {
		slog.Error("[ui] failed to get all tags", slog.Any("err", err))
		views.NewTagList(s.baseView, nil, err).Render(rw, r)
		return
	}
	views.NewTagList(s.baseView, tags, nil).Render(rw, r)
}

func (s *Handler) tagSearch(rw http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	input := &dusk.SearchFilters{
		Search: readString(qs, "itemSearch", ""),
	}

	tags, err := s.db.GetAllTags(input)
	if err != nil {
		log.Println(err)
		views.TagSearchResults(nil, err).Render(r.Context(), rw)
		return
	}
	views.TagSearchResults(tags, nil).Render(r.Context(), rw)
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
	views.NewTag(s.baseView, tag, &dusk.BooksPage{
		Page: dusk.Page{
			Total: int64(len(books)),
		},
		Books: books,
	}, nil).Render(rw, r)
}
