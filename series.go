package dusk

import (
	"fmt"

	"github.com/kencx/dusk/validator"
	"github.com/kennygrant/sanitize"
)

type Series struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Series) Slugify() string {
	return sanitize.Path(fmt.Sprintf("%s-%d", a.Name, a.Id))
}

func (t Series) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(t.Name != "", "name", "value is missing")
	return err
}
