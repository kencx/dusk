package dusk

import (
	"fmt"

	"github.com/kencx/dusk/validator"
	"github.com/kennygrant/sanitize"
)

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Authors []*Author

func (a Author) Slugify() string {
	return sanitize.Path(fmt.Sprintf("%s-%d", a.Name, a.ID))
}

func (a Author) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(a.Name != "", "name", "value is missing")
	return err
}
