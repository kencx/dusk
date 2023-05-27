package dusk

import "dusk/validator"

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Authors []*Author

func (a *Author) Validate(v *validator.Validator) {
	v.Check(a.Name != "", "name", "value is missing")
}
