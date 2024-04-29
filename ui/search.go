package ui

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"

	ol "github.com/kencx/dusk/integrations/openlibrary"
)

func (s *Handler) search(rw http.ResponseWriter, r *http.Request) {
	value := r.FormValue("openlibrary")

	ok, err := util.IsbnCheck(value)
	if err != nil {
		slog.Error("failed to fetch isbn", slog.Any("err", err))
		views.ImportResultsError(rw, r, err)
	}

	if ok {
		metadata, err := ol.FetchByIsbn(value)
		if err != nil {
			// TODO openlibrary request err
			slog.Error("failed to fetch isbn", slog.Any("err", err))
			views.ImportResultsError(rw, r, err)
			return
		}
		results := ol.QueryResults{metadata}
		slog.Info("Fetched results", slog.Any("results", results[0]))
		views.ImportResults(results, "", nil).Render(r.Context(), rw)

	} else {
		results, err := ol.FetchByQuery(value)
		if err != nil {
			slog.Error("failed to fetch query", slog.Any("err", err))
			views.ImportResultsError(rw, r, err)
			return
		}
		slog.Info("Fetched results", slog.Any("results", results))
		views.ImportResults(*results, "", nil).Render(r.Context(), rw)
	}
}

func (s *Handler) searchAddResult(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("result")

	// TODO We are fetching this endpoint and performing the same operations twice. It
	// would be good if we can cache the previously fetched data in importOpenLibrary on
	// the client side to send it here. This might require Alpine.js.

	// TODO handle book already added

	metadata, err := ol.FetchByIsbn(isbn)
	if err != nil {
		// TODO openlibrary request err
		views.ImportResultsError(rw, r, err)
		return
	}

	b := metadata.ToBook()

	errMap := validator.Validate(b)
	if len(errMap) > 0 {
		slog.Error("book validation failed", slog.Any("err", errMap))
		views.ImportResultsError(rw, r, errors.New("TODO"))
		return
	}

	book, err := s.db.CreateBook(b)
	if err != nil {
		slog.Error("create book failed", slog.Any("err", err))
		views.ImportResultsError(rw, r, err)
		return
	}

	// if b.Cover.Valid {
	// 	if err := s.fs.UploadBookCoverFromUrl(b.Cover.String, book); err != nil {
	// 		slog.Warn("failed to upload cover image", slog.Any("err", err))
	// 		views.ImportResultsError(rw, r, err)
	// 		return
	// 	}
	// }

	rawMessage := fmt.Sprintf("<p><a href=\"/book/%d\">%s</a> added</p>", book.ID, book.Title)
	views.ImportResultsMessage(rw, r, ol.QueryResults{metadata}, rawMessage)
}
