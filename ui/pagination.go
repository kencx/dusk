package ui

import (
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/page"
)

var (
	defaultAfterId  = 0
	defaultLimit    = 30
	defaultSort     = "name"
	defaultBookSort = "title"
)

func defaultFilters() *dusk.Filters {
	return &dusk.Filters{
		AfterId:      defaultAfterId,
		Limit:        defaultLimit,
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
			AfterId:      request.QueryInt(qs, page.After, defaultAfterId),
			Limit:        request.QueryInt(qs, page.Limit, defaultLimit),
			Sort:         request.QueryString(qs, page.Sort, defaultSort),
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
				AfterId:      request.QueryInt(qs, page.After, defaultAfterId),
				Limit:        request.QueryInt(qs, page.Limit, defaultLimit),
				Sort:         request.QueryString(qs, page.Sort, defaultBookSort),
				SortSafeList: dusk.DefaultSafeList(),
			},
		},
	}
}
