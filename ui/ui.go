package ui

import (
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/integration"
	"github.com/kencx/dusk/ui/views"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	GetBook(id int64) (*dusk.Book, error)
	GetAllBooks(filters *dusk.BookFilters) (*dusk.BooksPage, error)
	CreateBook(b *dusk.Book) (*dusk.Book, error)
	UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error)
	DeleteBook(id int64) error

	GetAuthor(id int64) (*dusk.Author, error)
	GetAllAuthors(filters *dusk.SearchFilters) (dusk.Authors, error)
	GetAllBooksFromAuthor(id int64) (dusk.Books, error)
	CreateAuthor(a *dusk.Author) (*dusk.Author, error)
	UpdateAuthor(id int64, a *dusk.Author) (*dusk.Author, error)
	DeleteAuthor(id int64) error

	GetTag(id int64) (*dusk.Tag, error)
	GetAllTags(filters *dusk.SearchFilters) (dusk.Tags, error)
	GetAllBooksFromTag(id int64) (dusk.Books, error)
	CreateTag(t *dusk.Tag) (*dusk.Tag, error)
	UpdateTag(id int64, t *dusk.Tag) (*dusk.Tag, error)
	DeleteTag(id int64) error
}

type Fetcher interface {
	FetchByIsbn(isbn string) (*integration.Metadata, error)
	FetchByQuery(query string) (*integration.QueryResults, error)
}

type Handler struct {
	db       Store
	fs       *file.Service
	f        Fetcher
	baseView views.BaseView
}

func Router(revision string, db Store, fs *file.Service, f Fetcher) chi.Router {
	bv := views.NewBaseView(revision)
	s := Handler{db, fs, f, bv}
	ui := chi.NewRouter()

	staticFiles(ui)
	dfs := http.FileServer(http.Dir(fs.Directory))
	ui.Handle("/files/*", http.StripPrefix("/files/", dfs))

	ui.HandleFunc("/", s.index)
	ui.Route("/b", func(c chi.Router) {
		c.Get("/{slug:[a-zA-Z0-9-]+}", s.bookPage)
		// c.Put("/{slug:[a-zA-Z0-9-]+}", s.updateBook)
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
		// c.Put("/{slug:[a-zA-Z0-9-]+}", s.updatetag)
		// c.Delete("/{slug:[a-zA-Z0-9-]+}", s.deletetag)
		c.Get("/search", s.tagSearch)
	})

	ui.HandleFunc("/import", s.importIndex)

	ui.Route("/search", func(c chi.Router) {
		c.Get("/", s.searchPage)
		c.Post("/", s.search)
		// TODO with pagination
		// c.Get("/results/{page:[0-9]+}", s.searchResultsPage)
		c.Post("/add", s.searchAddResult)
	})

	ui.Route("/upload", func(c chi.Router) {
		c.Get("/", s.uploadPage)
		c.Post("/upload", s.upload)
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

func (s *Handler) notFound(rw http.ResponseWriter, r *http.Request) {
	views.NotFound().Render(r.Context(), rw)
}
