package ui

import (
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/ui/views"

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
	GetAllBooksFromTag(id int64) (dusk.Books, error)
	CreateTag(t *dusk.Tag) (*dusk.Tag, error)
	UpdateTag(id int64, t *dusk.Tag) (*dusk.Tag, error)
	DeleteTag(id int64) error
}

type Handler struct {
	db Store
	fs *file.Service
}

func Router(db Store, fs *file.Service) chi.Router {
	s := Handler{db, fs}
	ui := chi.NewRouter()

	staticFiles(ui)
	dfs := http.FileServer(http.Dir(fs.Directory))
	ui.Handle("/files/*", http.StripPrefix("/files/", dfs))

	ui.HandleFunc("/", s.index)
	ui.Route("/b", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.bookPage)
		// c.Post("/{id[0-9]+}", s.formAddBook)
		// c.Put("/{id[0-9]+}", s.formUpdateBook)
		c.Delete("/{id:[0-9]+}", s.deleteBook)
	})

	ui.Route("/import", func(c chi.Router) {
		c.Get("/", s.importTabsPage)
		// c.Post("/goodreads", s.importGoodreads)
		// c.Post("/calibre", s.importCalibre)
	})
	ui.Route("/search", func(c chi.Router) {
		c.Post("/", s.search)
		c.Post("/add", s.searchAddResult)
	})
	ui.Post("/upload", s.upload)

	ui.HandleFunc("/authors", s.authorList)
	ui.Route("/a", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.authorPage)
		// c.Post("/{id[0-9]+}", s.formAddAuthor)
		// c.Put("/{id[0-9]+}", s.formUpdateAuthor)
		// c.Delete("/{id[0-9]+}", s.formDeleteAuthor)
	})

	ui.HandleFunc("/tags", s.tagList)
	ui.Route("/t", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.tagPage)
		// c.Post("/{id[0-9]+}", s.formAddtag)
		// c.Put("/{id[0-9]+}", s.formUpdatetag)
		// c.Delete("/{id[0-9]+}", s.formDeletetag)
	})

	ui.NotFound(s.notFound)
	return ui
}

func (s *Handler) notFound(rw http.ResponseWriter, r *http.Request) {
	views.NotFound().Render(r.Context(), rw)
}
