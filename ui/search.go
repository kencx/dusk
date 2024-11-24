package ui

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/integration"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) searchPage(rw http.ResponseWriter, r *http.Request) {
	views.NewImportIndex(s.base, "search", nil).Render(rw, r)
}

// TODO handle timeouts, 5XX errors
func (s *Handler) search(rw http.ResponseWriter, r *http.Request) {
	value := r.FormValue("search")

	isbnValid, err := util.IsbnCheck(value)
	if err != nil {
		slog.Error("[search] invalid isbn", slog.String("isbn", value), slog.Any("err", err))
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	if isbnValid {
		metadata, err := s.f.FetchByIsbn(value)
		if err != nil {
			slog.Error("[search] Failed to fetch by isbn", slog.String("isbn", value), slog.Any("err", err))
			views.SearchError(err).Render(r.Context(), rw)
			return
		}

		results := integration.QueryResults{metadata}
		views.SearchResults(results).Render(r.Context(), rw)

	} else {
		results, err := s.f.FetchByQuery(value)
		if err != nil {
			slog.Error("[search] Failed to fetch by query", slog.String("query", value), slog.Any("err", err))
			views.SearchError(err).Render(r.Context(), rw)
			return
		}

		slog.Debug("[search] Search successful", slog.Int("total", len(*results)), slog.String("query", value), slog.String("fetcher", ""))
		views.SearchResults(*results).Render(r.Context(), rw)
	}
}

func (s *Handler) searchAddResult(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("result")

	var readStatus dusk.ReadStatus
	switch r.FormValue("read-status") {
	case "unread":
		readStatus = dusk.Unread
	case "read":
		readStatus = dusk.Read
	case "reading":
		readStatus = dusk.Reading
	default:
		readStatus = dusk.Unread
	}

	// TODO We are fetching this endpoint and performing the same operations twice. It
	// would be good if we can cache the previously fetched data in importOpenLibrary on
	// the client side to send it here. This might require Alpine.js.

	metadata, err := s.f.FetchByIsbn(isbn)
	if err != nil {
		slog.Error(err.Error())
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	b := metadata.ToBook()
	b.DateAdded = null.TimeFrom(time.Now())
	b.Status = readStatus

	errMap := validator.Validate(b)
	if len(errMap) > 0 {
		slog.Error("failed to validate book", slog.Any("err", errMap))
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	book, err := s.db.CreateBook(b)
	if err != nil {
		if errors.Is(err, dusk.ErrIsbnExists) {
			slog.Error("failed to create book", slog.Any("err", err))
			SendToastMessage(rw, r, "Book already exists!")
			return
		}

		slog.Error("failed to create book", slog.Any("err", err))
		views.SearchError(err).Render(r.Context(), rw)
		return
	}

	if b.Cover.Valid {
		if err := s.fs.UploadCoverFromUrl(b.Cover.ValueOrZero(), book); err != nil {
			slog.Warn("failed to download cover", slog.Any("err", err))
			SendToastMessage(rw, r, "Failed to download cover!")
			return
		}

		// properly update cover filepath on db
		if _, err := s.db.UpdateBook(book.Id, book); err != nil {
			slog.Warn("failed to update book cover in database", slog.Any("err", err))
		}
	}

	rawMessage := fmt.Sprintf(`Book <a href="/b/%s">%s</a> added`, book.Slugify(), book.Title)
	SendToastRawMessage(rw, r, rawMessage)
}
