package dusk

import (
	"strings"

	"github.com/kencx/dusk/validator"
)

type Filters struct {
	AfterId      int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f *Filters) Valid() validator.ErrMap {
	errMap := validator.New()

	errMap.Check(f.AfterId >= 0, "after_id", "must be >= 0")
	errMap.Check(f.AfterId <= 10_000_000, "after_id", "must be <= 10 million")
	errMap.Check(f.PageSize > 0, "page_size", "must be > 0")
	errMap.Check(f.PageSize <= 1000, "page_size", "must be <= 1000")
	errMap.Check(validator.In(f.Sort, f.SortSafeList), "sort", "invalid sort value")

	return errMap
}

func (f Filters) SortColumn() string {
	for _, sv := range f.SortSafeList {
		if f.Sort == sv {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	// panic in case of SQL injection
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

type BookFilters struct {
	Title  string
	Author string
	Genre  string
	Tag    string
	Series string
	Filters
}
