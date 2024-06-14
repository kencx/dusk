package dusk

import (
	"net/url"
	"strconv"
)

type Pager interface {
	Next() bool
	Previous() bool
	IsFirst() bool
	IsLast() bool
	NumOfPages() int
}

type PageInfo struct {
	Limit       int
	TotalCount  int64
	FirstRowNo  int64
	LastRowNo   int64
	QueryParams url.Values
}

type Page[T any] struct {
	*PageInfo
	Items []T
}

func (p *Page[T]) Next() string {
	if p.IsLast() {
		return ""
	}

	if p.QueryParams.Has("after_id") {
		p.QueryParams.Set("after_id", strconv.Itoa(int(p.LastRowNo)))
	}

	return p.QueryParams.Encode()
}

func (p *Page[T]) Previous() string {
	if p.IsFirst() {
		return ""
	}

	if p.QueryParams.Has("after_id") {
		p.QueryParams.Set("after_id", strconv.Itoa(int(p.FirstRowNo)-p.Limit-1))
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
	return (int(p.TotalCount) / p.Limit) + 1
}
