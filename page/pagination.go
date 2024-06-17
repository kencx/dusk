package page

import (
	"net/url"
	"strconv"

	"github.com/kencx/dusk"
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

func New[T any](total, first, last int, filters *dusk.Filters, items []T) *Page[T] {
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
		p.QueryParams.Set(After, strconv.Itoa(int(p.FirstRowNo)-p.Limit-1))
	}
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
