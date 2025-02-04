package ui

import (
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/integration"
	"github.com/kencx/dusk/ui/shared"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db   dusk.Store
	fs   *file.Service
	f    integration.Fetchers
	base shared.Base
}

func Router(revision string, db dusk.Store, fs *file.Service, f integration.Fetchers) chi.Router {
	base := shared.NewBase(revision)
	s := Handler{db, fs, f, base}
	ui := chi.NewRouter()

	// middlewares
	ui.Use(response.NoCache)

	// static
	ui.Mount("/static", staticRouter(604800))
	ui.Mount("/files", filesRouter(fs, 604800))

	ui.HandleFunc("/", s.index)
	ui.Route("/b", func(c chi.Router) {
		c.Get("/{slug:[a-zA-Z0-9-]+}", s.bookPage)
		c.Get("/{slug:[a-zA-Z0-9-]+}/edit", s.editBookForm)
		c.Put("/{slug:[a-zA-Z0-9-]+}", s.updateBook)
		c.Put("/{slug:[a-zA-Z0-9-]+}/status", s.updateBookStatus)
		c.Delete("/{slug:[a-zA-Z0-9-]+}", s.deleteBook)
		c.Get("/search", s.bookSearch)

		// c.Get("/partials/rating", s.bookRatingPartial)
		// c.Get("/partials/tags", s.bookTagsPartial)
		// c.Get("/partials/cover", s.bookCoverPartial)
		// c.Get("/partials/actions", s.bookActionsPartial)
		// c.Get("/modal/delete", s.bookDeleteModal)
	})

	ui.HandleFunc("/authors", s.authorList)
	ui.Route("/a", func(c chi.Router) {
		c.Get("/{slug:[a-zA-Z0-9-]+}", s.authorPage)
		// c.Put("/{slug:[a-zA-Z0-9-]+}", s.updateAuthor)
		// c.Delete("/{slug:[a-zA-Z0-9-]+}", s.deleteAuthor)
		c.Get("/search", s.authorSearch)
	})

	ui.HandleFunc("/tags", s.tagList)
	ui.Route("/t", func(c chi.Router) {
		c.Get("/{slug:[a-zA-Z0-9-]+}", s.tagPage)
		c.Get("/all", s.tagDataList)
		// c.Put("/{slug:[a-zA-Z0-9-]+}", s.updatetag)
		// c.Delete("/{slug:[a-zA-Z0-9-]+}", s.deletetag)
		c.Get("/search", s.tagSearch)
	})

	ui.HandleFunc("/import", s.importIndex)

	ui.Route("/search", func(c chi.Router) {
		c.Get("/", s.searchPage)
		c.Get("/import", s.search)
		c.Post("/add", s.searchAddResult)
	})

	ui.Route("/upload", func(c chi.Router) {
		c.Get("/", s.uploadPage)
		c.Post("/", s.upload)
	})

	// ui.Route("/manual", func(c chi.Router) {
	// 	c.Get("/", s.manualPage)
	// 	c.Post("/manual", s.manual)
	// })

	ui.Route("/goodreads", func(c chi.Router) {
		c.Get("/", s.goodreadsPage)
		c.Post("/", s.goodreads)
	})

	// ui.Route("/calibre", func(c chi.Router) {
	// 	c.Get("/", s.calibrePage)
	// 	c.Post("/", s.calibre)
	// })

	ui.NotFound(s.notFound)
	return ui
}
