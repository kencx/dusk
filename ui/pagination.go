package ui

import (
	"net/http"

	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/page"
)

var (
	defaultAfterId  = 0
	defaultLimit    = 30
	defaultSort     = "name"
	defaultBookSort = "title"
)

func defaultFilters() *filters.Filters {
	return &filters.Filters{
		AfterId:      defaultAfterId,
		Limit:        defaultLimit,
		Sort:         defaultSort,
		SortSafeList: filters.DefaultSafeList(),
	}
}

func defaultSearchFilters() *filters.Search {
	return &filters.Search{
		Filters: *defaultFilters(),
	}
}

func defaultBookFilters() *filters.Book {
	bf := &filters.Book{
		Search: filters.Search{
			Filters: *defaultFilters(),
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
		Filters: filters.Filters{
			AfterId:      request.QueryInt(qs, page.After, defaultAfterId),
			Limit:        request.QueryInt(qs, page.Limit, defaultLimit),
			Sort:         request.QueryString(qs, page.Sort, defaultSort),
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
			Filters: filters.Filters{
				AfterId:      request.QueryInt(qs, page.After, defaultAfterId),
				Limit:        request.QueryInt(qs, page.Limit, defaultLimit),
				Sort:         request.QueryString(qs, page.Sort, defaultBookSort),
				SortSafeList: filters.DefaultSafeList(),
			},
		},
	}
}
