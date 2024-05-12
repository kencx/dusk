package ui

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kencx/dusk/integration"
	gb "github.com/kencx/dusk/integration/googlebooks"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) searchPage(rw http.ResponseWriter, r *http.Request) {
	views.ImportIndex("search").Render(r.Context(), rw)
}

// TODO handle timeouts, 5XX errors
func (s *Handler) search(rw http.ResponseWriter, r *http.Request) {
	value := r.FormValue("search")

	ok, err := util.IsbnCheck(value)
	if err != nil {
		slog.Error("[search] invalid isbn", slog.String("isbn", value))
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	if ok {
		metadata, err := gb.FetchByIsbn(value)
		if err != nil {
			slog.Error(err.Error())
			views.SearchError(err).Render(r.Context(), rw)
			return
		}

		results := integration.QueryResults{metadata}
		views.SearchResults(results).Render(r.Context(), rw)

	} else {
		results, err := gb.FetchByQuery(value)
		if err != nil {
			slog.Error(err.Error())
			views.SearchError(err).Render(r.Context(), rw)
			return
		}

		slog.Debug(fmt.Sprintf("Fetched %d results", len(*results)))
		views.SearchResults(*results).Render(r.Context(), rw)
	}
}

func (s *Handler) searchAddResult(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("result")
	readStatus := r.FormValue("read-status")

	// TODO We are fetching this endpoint and performing the same operations twice. It
	// would be good if we can cache the previously fetched data in importOpenLibrary on
	// the client side to send it here. This might require Alpine.js.

	metadata, err := gb.FetchByIsbn(isbn)
	if err != nil {
		slog.Error(err.Error())
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	b := metadata.ToBook()
	b.Tag = append(b.Tag, readStatus)

	errMap := validator.Validate(b)
	if len(errMap) > 0 {
		slog.Error("failed to validate book", slog.Any("err", errMap))
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	_, err = s.db.CreateBook(b)
	if err != nil {
		// TODO handle book that already exists
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			slog.Error("failed to create book", slog.Any("err", err))
			SendToastMessage(rw, r, "Book already exists!")
			return
		}

		slog.Error("failed to create book", slog.Any("err", err))
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	// TODO download book cover
	// if b.Cover.Valid {
	// 	if err := s.fs.UploadBookCoverFromUrl(b.Cover.String, book); err != nil {
	// 		slog.Warn("failed to upload cover image", slog.Any("err", err))
	// 		views.ImportResultsError(rw, r, err)
	// 		return
	// 	}
	// }

	SendToastMessage(rw, r, "Book added!")
}
