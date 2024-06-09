package dusk

type Page struct {
	Size       int
	Query      string
	Total      int64
	FirstRowNo int64
	LastRowNo  int64
}

func (p *Page) First() bool {
	return p.FirstRowNo <= 1
}

func (p *Page) Last() bool {
	return p.LastRowNo >= p.Total
}

type BooksPage struct {
	Page
	Books Books
}

type AuthorsPage struct {
	Page
	Authors Authors
}

type TagsPage struct {
	Page
	Tags Tags
}

type ItemsPage interface {
	BooksPage | AuthorsPage | TagsPage
}
