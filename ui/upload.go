package ui

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) upload(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "upload", "application/")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	b, err := s.fs.UploadBook(f)
	if err != nil {
		slog.Error("[UI] Failed to upload file", slog.Any("err", err))
		views.ImportResultsError(rw, r, err)
		return
	}

	res, err := s.db.CreateBook(b)
	if err != nil {
		slog.Error("[UI] Failed to create book", slog.Any("err", err))
		views.ImportResultsError(rw, r, err)
		return
	}

	rawMessage := fmt.Sprintf("<p><a href=\"/book/%d\">%s</a> added</p>", res.ID, res.Title)
	views.ImportResultsMessage(rw, r, nil, rawMessage)
}
