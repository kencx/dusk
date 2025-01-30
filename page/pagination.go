package page

import (
	"net/url"
	"strconv"

	"github.com/kencx/dusk/filters"
)

const (
	// query parameters
	After = "after"
	Limit = "limit"
	Sort  = "sort"
)

type Pager interface {
	Next() bool
	Previous() bool
	IsFirst() bool
	IsLast() bool
	NumOfPages() int
}

type Info struct {
	Limit       int
	TotalCount  int
	FirstRowNo  int
	LastRowNo   int
	QueryParams url.Values
}

type Page[T any] struct {
	*Info
	Items []T
}

func New[T any](total, first, last int, filters *filters.Base, items []T) *Page[T] {
	qp := make(url.Values)
	qp.Add(After, strconv.Itoa(filters.AfterId))
	qp.Add(Limit, strconv.Itoa(filters.Limit))
	qp.Add(Sort, filters.Sort)

	return &Page[T]{
		Info: &Info{
			Limit:       min(total, filters.Limit),
			TotalCount:  total,
			FirstRowNo:  first,
			LastRowNo:   last,
			QueryParams: qp,
		},
		Items: items,
	}
}

func Single[T any](filters *filters.Base, item T) *Page[T] {
	return &Page[T]{
		Info: &Info{
			Limit:       1,
			TotalCount:  1,
			FirstRowNo:  1,
			LastRowNo:   1,
			QueryParams: nil,
		},
		Items: []T{item},
	}
}

func NewEmpty[T any]() *Page[T] {
	return &Page[T]{}
}

func (p *Page[T]) Empty() bool {
	if p.Info != nil {
		return len(p.Items) == 0 && p.TotalCount == 0
	}
	return len(p.Items) == 0
}

func (p *Page[T]) Next() string {
	if p.IsLast() {
		return ""
	}

	if p.QueryParams.Has(After) {
		p.QueryParams.Set(After, strconv.Itoa(int(p.LastRowNo)))
	}
	return p.QueryParams.Encode()
}

func (p *Page[T]) Previous() string {
	if p.IsFirst() {
		return ""
	}

	if p.QueryParams.Has(After) {
		p.QueryParams.Set(After, strconv.Itoa(max(0, int(p.FirstRowNo)-p.Limit-1)))
	}
	return p.QueryParams.Encode()
}

func (p *Page[T]) First() string {
	p.QueryParams.Set(After, "0")
	return p.QueryParams.Encode()
}

func (p *Page[T]) Last() string {
	p.QueryParams.Set(After, strconv.Itoa((p.TotalCount/p.Limit)*p.Limit))
	return p.QueryParams.Encode()
}

func (p Page[T]) IsFirst() bool {
	return p.FirstRowNo <= 1
}

func (p Page[T]) IsLast() bool {
	return p.LastRowNo >= p.TotalCount
}

func (p Page[T]) NumOfPages() int {
	return (p.TotalCount / p.Limit) + 1
}
