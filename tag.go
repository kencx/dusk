package dusk

import (
	"github.com/kencx/dusk/validator"
)

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Tags []*Tag

func (t Tag) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(t.Name != "", "name", "value is missing")
	return err
}
