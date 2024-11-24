package ui

import (
	"net/http"

	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/page"
)

var (
	defaultFilters = &filters.Base{
		AfterId:      0,
		Limit:        30,
		Sort:         "name",
		SortSafeList: filters.DefaultSafeList(),
	}
	defaultBookSort = "title"
)

func defaultSearchFilters() *filters.Search {
	return &filters.Search{
		Base: *defaultFilters,
	}
}

func defaultBookFilters() *filters.Book {
	bf := &filters.Book{
		Search: filters.Search{
			Base: *defaultFilters,
		},
	}
	bf.Sort = defaultBookSort
	return bf
}

func initSearchFilters(r *http.Request) *filters.Search {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	return &filters.Search{
		Search: request.QueryString(qs, "q", ""),
		Base: filters.Base{
			AfterId:      request.QueryInt(qs, page.After, defaultFilters.AfterId),
			Limit:        request.QueryInt(qs, page.Limit, defaultFilters.Limit),
			Sort:         request.QueryString(qs, page.Sort, defaultFilters.Sort),
			SortSafeList: filters.DefaultSafeList(),
		},
	}
}

func initBookFilters(r *http.Request) *filters.Book {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	return &filters.Book{
		Title:  request.QueryString(qs, "title", ""),
		Author: request.QueryString(qs, "author", ""),
		Tag:    request.QueryString(qs, "tag", ""),
		Series: request.QueryString(qs, "series", ""),
		Search: filters.Search{
			Search: request.QueryString(qs, "q", ""),
			Base: filters.Base{
				AfterId:      request.QueryInt(qs, page.After, defaultFilters.AfterId),
				Limit:        request.QueryInt(qs, page.Limit, defaultFilters.Limit),
				Sort:         request.QueryString(qs, page.Sort, defaultBookSort),
				SortSafeList: filters.DefaultSafeList(),
			},
		},
	}
}
