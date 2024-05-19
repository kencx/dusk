package dusk

import (
	"fmt"
	"strings"

	"github.com/kencx/dusk/validator"
	"github.com/kennygrant/sanitize"
)

type Author struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Authors []*Author

func (a Author) Slugify() string {
	name := strings.ReplaceAll(a.Name, ".", "")
	return sanitize.Path(fmt.Sprintf("%s-%d", name, a.Id))
}

func (a Author) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(a.Name != "", "name", "value is missing")
	return err
}
