package util

import (
	"strings"

	"github.com/guregu/null/v5"
)

type NullString null.String

func (n NullString) Split(sep string) []string {
	if n.Valid {
		return strings.Split(n.String, sep)
	}
	return nil
}
