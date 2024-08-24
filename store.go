package dusk

import (
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/page"
)

type Store interface {
	GetBook(id int64) (*Book, error)
	GetAllBooks(filters *filters.Book) (*page.Page[Book], error)
	CreateBook(b *Book) (*Book, error)
	UpdateBook(id int64, b *Book) (*Book, error)
	DeleteBook(id int64) error

	GetAuthor(id int64) (*Author, error)
	GetAuthorsFromBook(id int64) ([]Author, error)
	GetAllAuthors(filters *filters.Search) (*page.Page[Author], error)
	GetAllBooksFromAuthor(id int64, filters *filters.Book) (*page.Page[Book], error)
	CreateAuthor(a *Author) (*Author, error)
	UpdateAuthor(id int64, a *Author) (*Author, error)
	DeleteAuthor(id int64) error

	GetTag(id int64) (*Tag, error)
	GetTagsFromBook(id int64) ([]Tag, error)
	GetAllTags(filters *filters.Search) (*page.Page[Tag], error)
	GetAllBooksFromTag(id int64, filters *filters.Book) (*page.Page[Book], error)
	CreateTag(t *Tag) (*Tag, error)
	UpdateTag(id int64, t *Tag) (*Tag, error)
	DeleteTag(id int64) error
}
