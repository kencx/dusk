package dusk

import "dusk/validator"

type Tag struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DateAdded   string `json:"-"`
	DateUpdated string `json:"-"`
}

type Tags []*Tag

func (t *Tag) Validate(v *validator.Validator) {
	v.Check(t.Name != "", "name", "value is missing")
}
