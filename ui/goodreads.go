package ui

import (
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) goodreadsPage(rw http.ResponseWriter, r *http.Request) {
	views.ImportIndex("goodreads").Render(r.Context(), rw)
}

func (s *Handler) goodreads(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "goodreads", "text/")
	if err != nil {
		slog.Error("[goodreads] failed to import csv", slog.Any("err", err))
		views.GoodreadsError(err).Render(r.Context(), rw)
		return
	}

	books, err := s.fs.ReadGoodreadsCSV(f)
	if err != nil {
		slog.Error("[goodreads] failed to read csv", slog.Any("err", err))
		views.GoodreadsError(err).Render(r.Context(), rw)
		return
	}

	var errMap = make(map[int]error)

	// TODO run as background job OR show load bar and prevent user from navigating away

	// TODO when re-importing csvs, books without any isbn will not fail the isbn
	// constraint requirement and be imported twice
	for i, book := range books {
		_, err := s.db.CreateBook(book)
		if err != nil {
			slog.Error("[goodreads] failed to create book", slog.Any("err", err))
			errMap[i] = err
		}
	}

	// TODO download book covers
	views.GoodreadsResults(books, errMap).Render(r.Context(), rw)
}
