package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tdewolff/minify/v2"
)

//go:embed static/*
var staticFS embed.FS

func staticFiles(router *chi.Mux) {
	staticFs, err := fs.Sub(staticFS, "static")
	if err != nil {
		slog.Error("[UI] failed to locate \"static\" directory in embed FS")
		return
	}

	m := minify.New()
	sfs := http.FileServer(http.FS(staticFs))
	minified := m.Middleware(http.StripPrefix("/static/", sfs))
	router.Handle("/static/*", minified)
}
