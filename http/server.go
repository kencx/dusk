package http

import (
	"context"
	"dusk"
	"dusk/storage"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	idleTimeout      = 60 * time.Second
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
	router   *mux.Router
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func New(db *storage.Store) *Server {
	r := mux.NewRouter()
	s := &Server{
		Server: &http.Server{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
			Handler:      r,
		},
		router:   r,
		db:       db,
		InfoLog:  log.New(os.Stdout, "INFO ", log.LstdFlags),
		ErrorLog: log.New(os.Stderr, "ERROR ", log.LstdFlags),
	}
	s.RegisterRoutes()
	return s
}

func (s *Server) Run(port string) error {
	s.Addr = port
	err := s.ListenAndServe()
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
	router := s.router

	api := router.PathPrefix("/api").Subrouter()
	_ = api.PathPrefix("/books/").Subrouter()
}
