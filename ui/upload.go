package ui

import (
	"log/slog"
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) uploadPage(rw http.ResponseWriter, r *http.Request) {
	views.NewImportIndex(s.base, "upload", nil).Render(rw, r)
}

func (s *Handler) upload(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "upload", "application/")
	if err != nil {
		slog.Error("failed to read file", slog.Any("err", err))
		views.UploadError(err).Render(r.Context(), rw)
		return
	}

	b, err := s.fs.UploadNewBook(f)
	if err != nil {
		slog.Error("[UI] Failed to upload file", slog.Any("err", err))
		views.UploadError(err).Render(r.Context(), rw)
		return
	}

	res, err := s.db.CreateBook(b)
	if err != nil {
		slog.Error("[UI] Failed to create book", slog.Any("err", err))
		views.UploadError(err).Render(r.Context(), rw)
		return
	}

	if r.FormValue("multiple") == "on" {
		views.UploadSuccess(res).Render(r.Context(), rw)
		return
	}
	response.HxRedirect(rw, r, "/b/"+res.Slugify())
}
