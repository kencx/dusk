package ui

import (
	"dusk"
	"dusk/file"
	"dusk/ui/views"
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TODO minify static files

//go:embed static/*.js static/*.css
var staticFs embed.FS

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
	fw *file.Worker
}

func Router(db Store, fw *file.Worker) chi.Router {
	s := Handler{db, fw}
	ui := chi.NewRouter()

	serveStaticFiles(ui, fw.DataDir)

	ui.HandleFunc("/", s.index)
	ui.Route("/book", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.bookPage)
		// c.Post("/{id[0-9]+}", s.formAddBook)
		// c.Put("/{id[0-9]+}", s.formUpdateBook)
		c.Delete("/{id:[0-9]+}", s.deleteBook)
	})

	ui.Post("/upload", s.upload)
	ui.Route("/import", func(c chi.Router) {
		c.Get("/", s.importPage)
		c.Post("/openlibrary", s.importOpenLibrary)
		// c.Post("/goodreads", s.importGoodreads)
		// c.Post("/calibre", s.importCalibre)
		c.Post("/add", s.importAddResult)
	})

	ui.HandleFunc("/authors", s.authorList)
	ui.Route("/author", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.authorPage)
		// c.Post("/{id[0-9]+}", s.formAddAuthor)
		// c.Put("/{id[0-9]+}", s.formUpdateAuthor)
		// c.Delete("/{id[0-9]+}", s.formDeleteAuthor)
	})

	ui.HandleFunc("/tags", s.tagList)
	ui.Route("/tag", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.tagPage)
		// c.Post("/{id[0-9]+}", s.formAddtag)
		// c.Put("/{id[0-9]+}", s.formUpdatetag)
		// c.Delete("/{id[0-9]+}", s.formDeletetag)
	})

	ui.NotFound(s.notFound)
	return ui
}

func serveStaticFiles(router *chi.Mux, dataDir string) {
	fs := http.FileServerFS(staticFs)
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	dfs := http.FileServer(http.Dir(dataDir))
	router.Handle("/files/*", http.StripPrefix("/files/", dfs))
}

func (s *Handler) notFound(rw http.ResponseWriter, r *http.Request) {
	views.NotFound().Render(r.Context(), rw)
}
