package dusk

import (
	"testing"

	"github.com/kencx/dusk/null"
)

var (
	isbnPass   = null.StringFrom("0143039822")
	isbnFail   = null.StringFrom("abc")
	isbn13Pass = null.StringFrom("9780316129084")
)

func TestValidateBook(t *testing.T) {
	tests := []struct {
		name string
		book *Book
		err  map[string]string
	}{{
		name: "success",
		book: &Book{
			Title:  "FooBar",
			ISBN:   isbnPass,
			Author: []string{"John Doe"},
		},
		err: nil,
	}, {
		name: "no title",
		book: &Book{
			ISBN:   isbnPass,
			Author: []string{"John Doe"},
		},
		err: map[string]string{"title": "value is missing"},
	}, {
		name: "nil author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   isbnPass,
			Author: nil,
		},
		err: map[string]string{"author": "value is missing"},
	}, {
		name: "zero length author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   isbnPass,
			Author: []string{},
		},
		err: map[string]string{"author": "value is missing"},
	}, {
		name: "no isbn",
		book: &Book{
			Title:  "Foo Bar",
			Author: []string{"John Doe"},
		},
		err: map[string]string{"isbn or isbn13": "value is missing"},
	}, {
		name: "isbn regex fail",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   isbnFail,
			Author: []string{"John Doe"},
		},
		err: map[string]string{"isbn10": "invalid isbn"},
	}, {
		name: "multiple errors",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   isbnFail,
			Author: nil,
		},
		err: map[string]string{"author": "value is missing", "isbn10": "invalid isbn"},
	}, {
		name: "both isbn",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   isbnPass,
			ISBN13: isbn13Pass,
			Author: []string{"John Doe"},
		},
		err: nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMap := tt.book.Valid()

			if len(errMap) > 0 && tt.err == nil {
				t.Fatalf("expected no err, got %v", errMap)
			}

			if len(errMap) == 0 && tt.err != nil {
				t.Fatalf("expected err with %q, got nil", tt.err)
			}

			if len(errMap) > 0 && tt.err != nil {
				if len(errMap) != len(tt.err) {
					t.Fatalf("got %d errs, want %d errs", len(errMap), len(tt.err))
				}

				for k, v := range errMap {
					s, ok := tt.err[k]
					if !ok {
						t.Fatalf("err field missing %q", k)
					}

					if v != s {
						t.Fatalf("got %v, want %v error", v, s)
					}
				}
			}
		})
	}
}
