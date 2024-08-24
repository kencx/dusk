package api

import (
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db       dusk.Store
	fs       *file.Service
	revision string
}

func Router(revision string, db dusk.Store, fs *file.Service) chi.Router {
	s := Handler{db, fs, revision}
	api := chi.NewRouter()

	api.Route("/books", func(r chi.Router) {
		r.Get("/{id:[0-9]+}", s.GetBook)
		r.Get("/", s.GetAllBooks)
		r.Post("/", s.AddBook)
		r.Post("/{id:[0-9]+}/cover", s.AddBookCover)
		r.Post("/{id:[0-9]+}/format", s.AddBookFormat)
		r.Put("/{id:[0-9]+}", s.UpdateBook)
		r.Delete("/{id:[0-9]+}", s.DeleteBook)
	})

	api.Route("/authors", func(r chi.Router) {})
	api.Route("/tags", func(r chi.Router) {})
	return api
}
