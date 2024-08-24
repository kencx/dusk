package filters

import (
	"strings"

	"github.com/kencx/dusk/validator"
)

type Filters struct {
	AfterId      int
	Limit        int
	Sort         string
	SortSafeList []string
}

func DefaultSafeList() []string {
	return []string{"title", "-title", "name", "-name"}
}

func (f Filters) Valid() validator.ErrMap {
	errMap := validator.New()

	errMap.Check(f.AfterId >= 0, "after", "must be >= 0")
	errMap.Check(f.AfterId <= 10_000_000, "after", "must be <= 10 million")
	errMap.Check(f.Limit > 0, "limit", "must be > 0")
	errMap.Check(f.Limit <= 1000, "limit", "must be <= 1000")
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

func (f *Filters) Empty() bool {
	return f.Sort == ""
}
