package ui

import (
	"dusk"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	GetBook(id int64) (*dusk.Book, error)
	GetAllBooks() (dusk.Books, error)
	CreateBook(b *dusk.Book) (*dusk.Book, error)
	UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error)
	DeleteBook(id int64) error

	GetAuthor(id int64) (*dusk.Author, error)
	GetAllAuthors() (dusk.Authors, error)
	GetAllBooksFromAuthor(id int64) (dusk.Books, error)
	CreateAuthor(a *dusk.Author) (*dusk.Author, error)
	UpdateAuthor(id int64, a *dusk.Author) (*dusk.Author, error)
	DeleteAuthor(id int64) error

	GetTag(id int64) (*dusk.Tag, error)
	GetAllTags() (dusk.Tags, error)
	CreateTag(t *dusk.Tag) (*dusk.Tag, error)
	UpdateTag(id int64, t *dusk.Tag) (*dusk.Tag, error)
	DeleteTag(id int64) error
}

type Handler struct {
	db Store
}

func Routes(db Store) chi.Router {
	s := Handler{db: db}
	ui := chi.NewRouter()

	fs := http.FileServer(http.Dir("./ui/static"))
	ui.Handle("/static/*", http.StripPrefix("/static/", fs))

	ui.HandleFunc("/", s.indexView)
	ui.Route("/book", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.bookView)
		// c.Post("/{id[0-9]+}", s.formAddBook)
		// c.Put("/{id[0-9]+}", s.formUpdateBook)
		// c.Delete("/{id[0-9]+}", s.formDeleteBook)
	})

	ui.Route("/import", func(c chi.Router) {
		c.Get("/", s.importView)
		c.Post("/openlibrary", s.importOpenLibrary)
		// c.Post("/goodreads", s.importGoodreads)
		// c.Post("/calibre", s.importCalibre)
	})

	ui.Route("/author", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.authorView)
		// c.Post("/{id[0-9]+}", s.formAddAuthor)
		// c.Put("/{id[0-9]+}", s.formUpdateAuthor)
		// c.Delete("/{id[0-9]+}", s.formDeleteAuthor)
	})

	return ui
}
