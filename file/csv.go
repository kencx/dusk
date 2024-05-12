package file

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/integrations/goodreads"
)

func (w *Service) ReadGoodreadsCSV(payload *Payload) (dusk.Books, error) {
	cr := csv.NewReader(payload.File)

	// read header
	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("failed to parse csv: %w", err)
	}

	var books dusk.Books
	var failed int
	for {
		record, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		book, err := goodreads.RecordToBook(record)
		if err != nil {
			slog.Warn(
				"failed to convert record to book",
				slog.String("title", record[1]),
				slog.Any("err", err),
			)
			failed += 1
			continue
		}
		books = append(books, book)
	}

	slog.Info(
		"[csv] read csv completed",
		slog.Int("success", len(books)-failed),
		slog.Int("failed", failed),
		slog.String("filename", payload.Filename),
	)
	return books, nil
}
