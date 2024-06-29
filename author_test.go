package dusk

import (
	"testing"
)

func TestValidateAuthor(t *testing.T) {
	tests := []struct {
		name   string
		author *Author
		err    map[string]string
	}{{
		name: "success",
		author: &Author{
			Name: "John Doe",
		},
		err: nil,
	}, {
		name: "no name",
		author: &Author{
			Name: "",
		},
		err: map[string]string{"name": "value is missing"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMap := tt.author.Valid()

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

func TestEqual(t *testing.T) {
	tests := []struct {
		name   string
		author *Author
		other  *Author
		result bool
	}{{
		name: "match",
		author: &Author{
			Name: "John Doe",
		},
		other: &Author{
			Name: "Doe, John",
		},
		result: true,
	}, {
		name: "no match",
		author: &Author{
			Name: "Jane Adams",
		},
		other: &Author{
			Name: "Adams, John",
		},
		result: false,
	}, {
		name: "same name",
		author: &Author{
			Name: "Jane Adams",
		},
		other: &Author{
			Name: "Jane Adams",
		},
		result: true,
	}, {
		name: "same name reversed",
		author: &Author{
			Name: "Doe, Jane",
		},
		other: &Author{
			Name: "Doe, Jane",
		},
		result: true,
	}, {
		name: "match middle name",
		author: &Author{
			Name: "George R.R. Martin",
		},
		other: &Author{
			Name: "Martin, George R.R.",
		},
		result: true,
	}, {
		name: "match middle name second",
		author: &Author{
			Name: "George R.R. Martin",
		},
		other: &Author{
			Name: "R.R. Martin, George",
		},
		result: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.author.Equal(*tt.other) != tt.result {
				t.Fatalf("not match")
			}
		})
	}
}
