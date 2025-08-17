package state

import (
	"net/http"

	"github.com/kencx/dusk/validator"
)

type Book struct {
	Base
}

func NewBookState(r *http.Request) (*Base, validator.ErrMap) {
	view := r.URL.Query().Get("view")

	filters := getBookFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		return nil, errMap
	}

	return NewBase(view, filters), nil
}
