package http

import (
	"context"
	"dusk"
	"dusk/storage"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	idleTimeout      = time.Minute
	readWriteTimeout = 3 * time.Second
	closeTimeout     = 5 * time.Second
)

type Store interface {
	GetBook(id int64) (*dusk.Book, error)
	GetAllBooks() (dusk.Books, error)
	CreateBook(b *dusk.Book) (*dusk.Book, error)
	UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error)
	DeleteBook(id int64) error

	GetAuthor(id int64) (*dusk.Author, error)
	GetAllAuthors() (dusk.Authors, error)
	CreateAuthor(a *dusk.Author) (*dusk.Author, error)
	UpdateAuthor(id int64, a *dusk.Author) (*dusk.Author, error)
	DeleteAuthor(id int64) error

	GetTag(id int64) (*dusk.Tag, error)
	GetAllTags() (dusk.Tags, error)
	CreateTag(t *dusk.Tag) (*dusk.Tag, error)
	UpdateTag(id int64, t *dusk.Tag) (*dusk.Tag, error)
	DeleteTag(id int64) error
}

type Server struct {
	*http.Server
	db       Store
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func New(db *storage.Store) *Server {
	s := &Server{
		Server: &http.Server{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
			Handler:      chi.NewRouter(),
		},
		db:       db,
		InfoLog:  log.New(os.Stdout, "INFO ", log.LstdFlags),
		ErrorLog: log.New(os.Stderr, "ERROR ", log.LstdFlags),
	}
	s.RegisterRoutes()
	return s
}

func (s *Server) Run(port, tlsCert, tlsKey string) error {
	s.Addr = port

	var err error
	if tlsCert != "" && tlsKey != "" {
		err = s.ListenAndServeTLS(tlsCert, tlsKey)
	} else {
		err = s.ListenAndServe()
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	tc, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	return s.Shutdown(tc)
}

func (s *Server) RegisterRoutes() {
	r := s.Handler.(*chi.Mux)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	fs := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.HandleFunc("/", s.indexView)

	r.Route("/import", func(c chi.Router) {
		c.Get("/", s.importView)
		c.Post("/openlibrary", s.importOpenLibrary)
		// c.Post("/goodreads", s.importGoodreads)
		// c.Post("/calibre", s.importCalibre)
	})

	r.Route("/book", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.bookView)
		// c.Post("/{id[0-9]+}", s.formAddBook)
		// c.Put("/{id[0-9]+}", s.formUpdateBook)
		// c.Delete("/{id[0-9]+}", s.formDeleteBook)
	})

	r.Route("/author", func(c chi.Router) {
		c.Get("/{id:[0-9]+}", s.authorView)
		// c.Post("/{id[0-9]+}", s.formAddAuthor)
		// c.Put("/{id[0-9]+}", s.formUpdateAuthor)
		// c.Delete("/{id[0-9]+}", s.formDeleteAuthor)
	})

	api := chi.NewRouter()
	r.Mount("/api", api)

	api.Route("/books", func(r chi.Router) {
		r.Get("/{id:[0-9]+}", s.GetBook)
		r.Get("/", s.GetAllBooks)
		r.Post("/", s.AddBook)
		r.Put("/{id:[0-9]+}", s.UpdateBook)
		r.Delete("/{id:[0-9]+}", s.DeleteBook)
	})

	api.Route("/authors", func(r chi.Router) {})
	api.Route("/tags", func(r chi.Router) {})
}
