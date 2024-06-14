package ui

import (
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
)

var (
	defaultAfterId  = 0
	defaultPageSize = 30
	defaultSort     = "name"
	defaultBookSort = "title"
)

func defaultFilters() *dusk.Filters {
	return &dusk.Filters{
		AfterId:      defaultAfterId,
		PageSize:     defaultPageSize,
		Sort:         defaultSort,
		SortSafeList: dusk.DefaultSafeList(),
	}
}

func defaultSearchFilters() *dusk.SearchFilters {
	return &dusk.SearchFilters{
		Filters: *defaultFilters(),
	}
}

func defaultBookFilters() *dusk.BookFilters {
	bf := &dusk.BookFilters{
		SearchFilters: dusk.SearchFilters{
			Filters: *defaultFilters(),
		},
	}
	bf.Sort = defaultBookSort
	return bf
}

func initSearchFilters(r *http.Request) *dusk.SearchFilters {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	return &dusk.SearchFilters{
		Search: request.QueryString(qs, "q", ""),
		Filters: dusk.Filters{
			AfterId:      request.QueryInt(qs, "after_id", defaultAfterId),
			PageSize:     request.QueryInt(qs, "page_size", defaultPageSize),
			Sort:         request.QueryString(qs, "sort", defaultSort),
			SortSafeList: dusk.DefaultSafeList(),
		},
	}
}

func initBookFilters(r *http.Request) *dusk.BookFilters {
	qs := r.URL.Query()

	// TODO trim, escape and filter special chars
	return &dusk.BookFilters{
		Title:  request.QueryString(qs, "title", ""),
		Author: request.QueryString(qs, "author", ""),
		Tag:    request.QueryString(qs, "tag", ""),
		Series: request.QueryString(qs, "series", ""),
		SearchFilters: dusk.SearchFilters{
			Search: request.QueryString(qs, "q", ""),
			Filters: dusk.Filters{
				AfterId:      request.QueryInt(qs, "after_id", defaultAfterId),
				PageSize:     request.QueryInt(qs, "page_size", defaultPageSize),
				Sort:         request.QueryString(qs, "sort", defaultBookSort),
				SortSafeList: dusk.DefaultSafeList(),
			},
		},
	}
}
