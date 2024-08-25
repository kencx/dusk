package filters

import (
	"strings"

	"github.com/kencx/dusk/validator"
)

type Base struct {
	AfterId      int
	Limit        int
	Sort         string
	SortSafeList []string
}

func DefaultSafeList() []string {
	return []string{"title", "-title", "name", "-name"}
}

func (b Base) Valid() validator.ErrMap {
	errMap := validator.New()

	errMap.Check(b.AfterId >= 0, "after", "must be >= 0")
	errMap.Check(b.AfterId <= 10_000_000, "after", "must be <= 10 million")
	errMap.Check(b.Limit > 0, "limit", "must be > 0")
	errMap.Check(b.Limit <= 1000, "limit", "must be <= 1000")
	errMap.Check(validator.In(b.Sort, b.SortSafeList), "sort", "invalid sort value")

	return errMap
}

func (b *Base) Empty() bool {
	return b.Sort == ""
}

func (b Base) SortColumn() string {
	for _, sv := range b.SortSafeList {
		if b.Sort == sv {
			return strings.TrimPrefix(b.Sort, "-")
		}
	}
	// panic in case of SQL injection
	panic("unsafe sort parameter: " + b.Sort)
}

func (b Base) SortDirection() string {
	if strings.HasPrefix(b.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
