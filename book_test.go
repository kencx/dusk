package dusk

import (
	"testing"

	"dusk/validator"
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
			ISBN:   "100",
			Author: []string{"John Doe"},
		},
		err: nil,
	}, {
		name: "no title",
		book: &Book{
			ISBN:   "100",
			Author: []string{"John Doe"},
		},
		err: map[string]string{"title": "value is missing"},
	}, {
		name: "nil author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "100",
			Author: nil,
		},
		err: map[string]string{"author": "value is missing"},
	}, {
		name: "zero length author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "100",
			Author: []string{},
		},
		err: map[string]string{"author": "value is missing"},
	}, {
		name: "no isbn",
		book: &Book{
			Title:  "Foo Bar",
			Author: []string{"John Doe"},
		},
		err: map[string]string{"isbn": "value is missing"},
	}, {
		name: "isbn regex fail",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "abc",
			Author: []string{"John Doe"},
		},
		err: map[string]string{"isbn": "incorrect format"},
	}, {
		name: "multiple errors",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "abc",
			Author: nil,
		},
		err: map[string]string{"author": "value is missing", "isbn": "incorrect format"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			tt.book.Validate(v)

			if !v.Valid() && tt.err == nil {
				t.Fatalf("expected no err, got %v", v.Errors)
			}

			if v.Valid() && tt.err != nil {
				t.Fatalf("expected err with %q, got nil", tt.err)
			}

			if !v.Valid() && tt.err != nil {
				if len(v.Errors) != len(tt.err) {
					t.Fatalf("got %d errs, want %d errs", len(v.Errors), len(tt.err))
				}

				for k, v := range v.Errors {
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
