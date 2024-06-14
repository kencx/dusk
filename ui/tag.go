package ui

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) tagList(rw http.ResponseWriter, r *http.Request) {
	tags, err := s.db.GetAllTags(defaultSearchFilters())
	if err != nil {
		slog.Error("[ui] failed to get all tags", slog.Any("err", err))
		views.NewTagList(s.base, dusk.Page[dusk.Tag]{}, err).Render(rw, r)
		return
	}
	views.NewTagList(s.base, *tags, nil).Render(rw, r)
}

func (s *Handler) tagSearch(rw http.ResponseWriter, r *http.Request) {

	filters := initSearchFilters(r)
	if errMap := validator.Validate(filters.Filters); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.TagSearchResults(dusk.Page[dusk.Tag]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	page, err := s.db.GetAllTags(filters)
	if err != nil {
		if err == dusk.ErrNoRows {
			views.TagSearchResults(dusk.Page[dusk.Tag]{}, err).Render(r.Context(), rw)
			return
		} else {
			slog.Error("failed to get all tags", slog.Any("err", err))
			views.TagSearchResults(dusk.Page[dusk.Tag]{}, err).Render(r.Context(), rw)
			return
		}
	}

	// only return page partial when not querying for first page
	if !page.IsFirst() {
		views.TagListPage(*page).Render(r.Context(), rw)
	} else {
		views.TagSearchResults(*page, nil).Render(r.Context(), rw)
	}
}

func (s *Handler) tagPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	filters := initBookFilters(r)
	if errMap := validator.Validate(filters.Filters); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.TagSearchResults(dusk.Page[dusk.Tag]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	tag, err := s.db.GetTag(id)
	if err != nil {
		slog.Error("[ui] failed to get tag", slog.Int64("id", id), slog.Any("err", err))
		views.NewTag(s.base, dusk.Tag{}, dusk.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromTag(tag.Id, filters)
	if err != nil {
		slog.Error("[ui] failed to get books from tag", slog.Int64("id", id), slog.Any("err", err))
		views.NewTag(s.base, dusk.Tag{}, dusk.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}
	views.NewTag(s.base, *tag, *books, nil).Render(rw, r)
}
