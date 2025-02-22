package ui

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) tagList(rw http.ResponseWriter, r *http.Request) {
	filters := initSearchFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.TagSearchResults(page.Page[dusk.Tag]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	tags, err := s.db.GetAllTags(filters)
	if err != nil {
		slog.Error("[ui] failed to get all tags", slog.Any("err", err))
		views.NewTagList(s.base, page.Page[dusk.Tag]{}, filters.Base, err).Render(rw, r)
		return
	}
	views.NewTagList(s.base, *tags, filters.Base, nil).Render(rw, r)
}

func (s *Handler) tagDataList(rw http.ResponseWriter, r *http.Request) {
	tags, err := s.db.GetAllTags(nil)
	if err != nil {
		slog.Error("[ui] failed to get all tags", slog.Any("err", err))
		views.TagDataList(page.Page[dusk.Tag]{}).Render(r.Context(), rw)
		return
	}

	views.TagDataList(*tags).Render(r.Context(), rw)
}

func (s *Handler) tagSearch(rw http.ResponseWriter, r *http.Request) {
	// If not htmx request, return the full page instead of partial.
	// Required to support hx-push-urls
	if request.IsHtmxRequest(r) {
		s.tagList(rw, r)
		return
	}

	filters := initSearchFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.TagSearchResults(page.Page[dusk.Tag]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	p, err := s.db.GetAllTags(filters)
	if err != nil {
		slog.Error("failed to get all tags", slog.Any("err", err))
		views.TagSearchResults(page.Page[dusk.Tag]{}, err).Render(r.Context(), rw)
		return
	}
	views.TagSearchResults(*p, nil).Render(r.Context(), rw)
}

func (s *Handler) tagPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	filters := initBookFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.TagSearchResults(page.Page[dusk.Tag]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	tag, err := s.db.GetTag(id)
	if err != nil {
		slog.Error("[ui] failed to get tag", slog.Int64("id", id), slog.Any("err", err))
		views.NewTag(s.base, dusk.Tag{}, page.Page[dusk.Book]{}, filters.Base, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromTag(tag.Id, filters)
	if err != nil {
		slog.Error("[ui] failed to get books from tag", slog.Int64("id", id), slog.Any("err", err))
		views.NewTag(s.base, dusk.Tag{}, page.Page[dusk.Book]{}, filters.Base, err).Render(rw, r)
		return
	}
	views.NewTag(s.base, *tag, *books, filters.Base, nil).Render(rw, r)
}
