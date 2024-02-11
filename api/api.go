package api

import (
	"dusk"
	"dusk/file"
	"dusk/http/response"
	"dusk/util"
	"net/http"
	"time"

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
	fw *file.Worker
}

func Router(db Store, fw *file.Worker) chi.Router {
	s := Handler{db, fw}
	api := chi.NewRouter()

	api.Get("/ping", s.Healthcheck)

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

func (s *Handler) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	res, err := util.ToJSON(response.Envelope{
		"timestamp": time.Now().Unix(),
		"message":   "pong",
	})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}
