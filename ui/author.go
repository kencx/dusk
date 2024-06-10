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

func (s *Handler) authorList(rw http.ResponseWriter, r *http.Request) {
	authors, err := s.db.GetAllAuthors(dusk.DefaultSearchFilters())
	if err != nil {
		slog.Error("[ui] failed to get all authors", slog.Any("err", err))
		views.NewAuthorList(s.base, dusk.Page[dusk.Author]{}, err).Render(rw, r)
		return
	}
	views.NewAuthorList(s.base, *authors, nil).Render(rw, r)
}

func (s *Handler) authorSearch(rw http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	filters := &dusk.SearchFilters{
		Search: readString(qs, "itemSearch", ""),
		Filters: dusk.Filters{
			AfterId:      readInt(qs, "after_id", 0),
			PageSize:     readInt(qs, "page_size", 30),
			Sort:         readString(qs, "sort", "name"),
			SortSafeList: dusk.DefaultSafeList(),
		},
	}

	if errMap := validator.Validate(filters.Filters); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.AuthorSearchResults(dusk.Page[dusk.Author]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	page, err := s.db.GetAllAuthors(filters)
	if err != nil {
		if err == dusk.ErrNoRows {
			views.AuthorSearchResults(dusk.Page[dusk.Author]{}, err).Render(r.Context(), rw)
			return
		} else {
			slog.Error("failed to get all authors", slog.Any("err", err))
			views.AuthorSearchResults(dusk.Page[dusk.Author]{}, err).Render(r.Context(), rw)
			return
		}
	}

	// only return page partial when not querying for first page
	if !page.First() {
		views.AuthorListPage(*page).Render(r.Context(), rw)
	} else {
		views.AuthorSearchResults(*page, nil).Render(r.Context(), rw)
	}
}

func (s *Handler) authorPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	qs := r.URL.Query()
	filters := &dusk.BookFilters{
		SearchFilters: dusk.SearchFilters{
			Search: readString(qs, "itemSearch", ""),
			Filters: dusk.Filters{
				AfterId:      readInt(qs, "after_id", 0),
				PageSize:     readInt(qs, "page_size", 30),
				Sort:         readString(qs, "sort", "title"),
				SortSafeList: dusk.DefaultSafeList(),
			},
		},
	}

	if errMap := validator.Validate(filters.Filters); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		views.AuthorSearchResults(dusk.Page[dusk.Author]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	author, err := s.db.GetAuthor(id)
	if err != nil {
		slog.Error("[ui] failed to get author", slog.Int64("id", id), slog.Any("err", err))
		views.NewAuthor(s.base, dusk.Author{}, dusk.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}

	books, err := s.db.GetAllBooksFromAuthor(author.Id, filters)
	if err != nil {
		slog.Error("[ui] failed to get books from author", slog.Int64("id", id), slog.Any("err", err))
		views.NewAuthor(s.base, dusk.Author{}, dusk.Page[dusk.Book]{}, err).Render(rw, r)
		return
	}
	views.NewAuthor(s.base, *author, *books, nil).Render(rw, r)
}
