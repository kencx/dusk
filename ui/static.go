package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kencx/dusk/http/response"
	"github.com/tdewolff/minify/v2"
)

//go:embed static/*
var staticFS embed.FS

func staticRouter(cacheDuration int) *chi.Mux {
	static := chi.NewRouter()

	// static middleware
	static.Use(
		response.ETag,
		response.SetCache(cacheDuration),
	)

	static.Handle("/*", minified())
	return static
}

func minified() http.Handler {
	staticFs, err := fs.Sub(staticFS, "static")
	if err != nil {
		slog.Error("[UI] failed to locate \"static\" directory in embed FS")
		return nil
	}

	m := minify.New()
	sfs := http.FileServer(http.FS(staticFs))
	mini := m.Middleware(http.StripPrefix("/static/", sfs))
	return mini
}
