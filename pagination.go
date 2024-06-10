package dusk

type Pager interface {
	First() bool
	Last() bool
	LastNo() int
	TotalNo() int
}

type Page[T any] struct {
	Size       int
	Query      string
	Total      int64
	FirstRowNo int64
	LastRowNo  int64
	Items      []T
}

func NewPage[T any](items []T) Page[T] {
	return Page[T]{Items: items}
}

func (p Page[T]) First() bool {
	return p.FirstRowNo <= 1
}

func (p Page[T]) Last() bool {
	return p.LastRowNo >= p.Total
}

func (p Page[T]) LastNo() int {
	return int(p.LastRowNo)
}

func (p Page[T]) TotalNo() int {
	return int(p.Total)
}
