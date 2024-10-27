package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/http/response"
)

func filesRouter(fs *file.Service, cacheDuration int) *chi.Mux {
	files := chi.NewRouter()

	// files middleware
	files.Use(
		response.ETag,
		response.SetCache(cacheDuration),
	)

	dfs := http.FileServer(http.Dir(fs.Directory))
	files.Handle("/*", http.StripPrefix("/files/", dfs))
	return files
}
