package dusk

import "dusk/validator"

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Authors []*Author

func (a Author) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(a.Name != "", "name", "value is missing")
	return err
}
