package ui

import (
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/integration/goodreads"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) goodreadsPage(rw http.ResponseWriter, r *http.Request) {
	views.NewImportIndex(s.base, "goodreads", nil).Render(rw, r)
}

func (s *Handler) goodreads(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "goodreads", "text/")
	if err != nil {
		slog.Error("[goodreads] failed to import csv", slog.Any("err", err))
		views.GoodreadsError(err).Render(r.Context(), rw)
		return
	}

	books, err := goodreads.ReadCSV(f)
	if err != nil {
		slog.Error("[goodreads] failed to read csv", slog.Any("err", err))
		views.GoodreadsError(err).Render(r.Context(), rw)
		return
	}

	// TODO run as background job OR show load bar and prevent user from navigating away
	// TODO download book covers
	var errMap = make(map[int]error)

	// handling duplicates
	// for _, book := range books {
	// TODO when re-importing csvs, books without any isbn will NOT fail the isbn
	// constraint requirement and be imported twice

	// _, err := s.db.GetBookEqual(book)
	// if err != nil {
	// 	slog.Error("[goodreads] failed to get duplicates", slog.Any("err", err))
	// }

	// if len(dups) > 0 { }
	// }

	for i, book := range books {
		b, err := s.db.CreateBook(book)
		if err != nil {
			slog.Error("[goodreads] failed to create book", slog.Any("err", err))
			errMap[i] = err
		}
		book.Id = b.Id
	}

	views.GoodreadsResults(books, errMap).Render(r.Context(), rw)
}
