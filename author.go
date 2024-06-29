package dusk

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kencx/dusk/validator"
	"github.com/kennygrant/sanitize"
)

type Author struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Author) Slugify() string {
	name := strings.ReplaceAll(a.Name, ".", "")
	return sanitize.Path(fmt.Sprintf("%s-%d", name, a.Id))
}

func (a Author) Valid() validator.ErrMap {
	err := validator.New()
	err.Check(a.Name != "", "name", "value is missing")
	return err
}

func (a Author) Equal(other Author) bool {
	if a.Name == other.Name {
		return true
	}

	var split []string
	split = strings.Split(other.Name, ",")
	if len(split) < 1 {
		split = strings.Split(a.Name, ",")
	}

	if len(split) < 1 {
		return a.Name == other.Name
	}
	for i, s := range split {
		split[i] = strings.TrimSpace(s)
	}

	slices.Reverse(split)
	return a.Name == strings.Join(split, " ")
}
