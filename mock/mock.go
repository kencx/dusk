package mock

import "github.com/kencx/dusk"

type Store struct {
	GetBookFn     func(id int64) (*dusk.Book, error)
	GetAllBooksFn func() (dusk.Books, error)
	CreateBookFn  func(b *dusk.Book) (*dusk.Book, error)
	UpdateBookFn  func(id int64, b *dusk.Book) (*dusk.Book, error)
	DeleteBookFn  func(id int64) error

	GetAuthorFn             func(id int64) (*dusk.Author, error)
	GetAllAuthorsFn         func() (dusk.Authors, error)
	GetAllBooksFromAuthorFn func(id int64) (dusk.Books, error)
	CreateAuthorFn          func(a *dusk.Author) (*dusk.Author, error)
	UpdateAuthorFn          func(id int64, a *dusk.Author) (*dusk.Author, error)
	DeleteAuthorFn          func(id int64) error

	GetTagFn             func(id int64) (*dusk.Tag, error)
	GetAllTagsFn         func() (dusk.Tags, error)
	GetAllBooksFromTagFn func(id int64) (dusk.Books, error)
	CreateTagFn          func(a *dusk.Tag) (*dusk.Tag, error)
	UpdateTagFn          func(id int64, a *dusk.Tag) (*dusk.Tag, error)
	DeleteTagFn          func(id int64) error
}

func (s *Store) GetBook(id int64) (*dusk.Book, error) {
	return s.GetBookFn(id)
}

func (s *Store) GetAllBooks() (dusk.Books, error) {
	return s.GetAllBooksFn()
}

func (s *Store) CreateBook(b *dusk.Book) (*dusk.Book, error) {
	return s.CreateBookFn(b)
}

func (s *Store) UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error) {
	return s.UpdateBookFn(id, b)
}

func (s *Store) DeleteBook(id int64) error {
	return s.DeleteBookFn(id)
}

func (s *Store) GetAuthor(id int64) (*dusk.Author, error) {
	return s.GetAuthorFn(id)
}

func (s *Store) GetAllAuthors() (dusk.Authors, error) {
	return s.GetAllAuthorsFn()
}

func (s *Store) GetAllBooksFromAuthor(id int64) (dusk.Books, error) {
	return s.GetAllBooksFromAuthorFn(id)
}

func (s *Store) CreateAuthor(b *dusk.Author) (*dusk.Author, error) {
	return s.CreateAuthorFn(b)
}

func (s *Store) UpdateAuthor(id int64, b *dusk.Author) (*dusk.Author, error) {
	return s.UpdateAuthorFn(id, b)
}

func (s *Store) DeleteAuthor(id int64) error {
	return s.DeleteAuthorFn(id)
}

func (s *Store) GetTag(id int64) (*dusk.Tag, error) {
	return s.GetTagFn(id)
}

func (s *Store) GetAllTags() (dusk.Tags, error) {
	return s.GetAllTagsFn()
}

func (s *Store) GetAllBooksFromTag(id int64) (dusk.Books, error) {
	return s.GetAllBooksFromTagFn(id)
}

func (s *Store) CreateTag(b *dusk.Tag) (*dusk.Tag, error) {
	return s.CreateTagFn(b)
}

func (s *Store) UpdateTag(id int64, b *dusk.Tag) (*dusk.Tag, error) {
	return s.UpdateTagFn(id, b)
}

func (s *Store) DeleteTag(id int64) error {
	return s.DeleteTagFn(id)
}
