package api

import (
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/page"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	GetBook(id int64) (*dusk.Book, error)
	GetAllBooks(filters *dusk.BookFilters) (*page.Page[dusk.Book], error)
	CreateBook(b *dusk.Book) (*dusk.Book, error)
	UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error)
	DeleteBook(id int64) error

	GetAuthor(id int64) (*dusk.Author, error)
	GetAuthorsFromBook(id int64) ([]dusk.Author, error)
	GetAllAuthors(filters *dusk.SearchFilters) (*page.Page[dusk.Author], error)
	GetAllBooksFromAuthor(id int64, filters *dusk.BookFilters) (*page.Page[dusk.Book], error)
	CreateAuthor(a *dusk.Author) (*dusk.Author, error)
	UpdateAuthor(id int64, a *dusk.Author) (*dusk.Author, error)
	DeleteAuthor(id int64) error

	GetTag(id int64) (*dusk.Tag, error)
	GetTagsFromBook(id int64) ([]dusk.Tag, error)
	GetAllTags(filters *dusk.SearchFilters) (*page.Page[dusk.Tag], error)
	GetAllBooksFromTag(id int64, filters *dusk.BookFilters) (*page.Page[dusk.Book], error)
	CreateTag(t *dusk.Tag) (*dusk.Tag, error)
	UpdateTag(id int64, t *dusk.Tag) (*dusk.Tag, error)
	DeleteTag(id int64) error
}

type Handler struct {
	db       Store
	fs       *file.Service
	revision string
}

func Router(revision string, db Store, fs *file.Service) chi.Router {
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
