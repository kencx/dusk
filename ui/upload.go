package ui

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) uploadPage(rw http.ResponseWriter, r *http.Request) {
	views.Search{DefaultTab: "upload"}.Render(rw, r)
}

func (s *Handler) upload(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "upload", "application/")
	if err != nil {
		slog.Error("failed to read file", slog.Any("err", err))
		// views.ImportResultsError(rw, r, err)
		return
	}

	b, err := s.fs.UploadBook(f)
	if err != nil {
		slog.Error("[UI] Failed to upload file", slog.Any("err", err))
		// views.ImportResultsError(rw, r, err)
		return
	}

	res, err := s.db.CreateBook(b)
	if err != nil {
		slog.Error("[UI] Failed to create book", slog.Any("err", err))
		// views.ImportResultsError(rw, r, err)
		return
	}

	_ = fmt.Sprintf("<a href=\"/b/%s\">%s</a> added", res.Slugify(), res.Title)
	// views.ImportResultsMessage(rw, r, nil, rawMessage)
}
