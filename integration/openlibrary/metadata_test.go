package openlibrary

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestUnmarshalInvalidMetadata(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{{
		name:  "no author",
		input: []byte(`{"title": "Foo Bar","isbn_10": ["0123456789"]}`),
		err:   ErrInvalidResult,
	}, {
		name:  "no title",
		input: []byte(`{"authors": [{"Key": "John Adams"}], "isbn_10": ["0123456789"]}`),
		err:   ErrInvalidResult,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			var got OlMetadata

			err := json.Unmarshal(tt.input, &got)
			is.True(errors.Is(err, ErrInvalidResult))
		})
	}
}
