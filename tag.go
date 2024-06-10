package dusk

import (
	"fmt"
	"strings"

	"github.com/kencx/dusk/validator"
	"github.com/kennygrant/sanitize"
)

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Tag) Slugify() string {
	return sanitize.Path(fmt.Sprintf("%s-%d", a.Name, a.Id))
}

func (t Tag) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(t.Name != "", "name", "value is missing")
	return err
}

func (t Tag) Parent() string {
	return strings.SplitN(t.Name, ".", 2)[0]
}
