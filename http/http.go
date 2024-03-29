package http

import (
	"context"
	"net/http"
	"time"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/api"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/storage"
	"github.com/kencx/dusk/ui"

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

type Server struct {
	*http.Server
	db Store
	fs *file.Service
}

func New(db *storage.Store, fs *file.Service) *Server {
	s := &Server{
		Server: &http.Server{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
			Handler:      chi.NewRouter(),
		},
		db: db,
		fs: fs,
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

	r.Mount("/api", api.Router(s.db, s.fs))
	r.Mount("/", ui.Router(s.db, s.fs))
}
