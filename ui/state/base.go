package state

import (
	"github.com/kencx/dusk/filters"
)

type Base struct {
	View    string
	Filters filters.Filters
}

func NewBase(view string, f filters.Filters) *Base {
	return &Base{view, f}
}
