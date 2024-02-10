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
