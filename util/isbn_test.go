package util

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestIsbnCheck(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   bool
		err    error
	}{{
		name: "isbn10",
		values: []string{
			"0143039822",
			"048624895X",
			"ISBN:0143039822",
			"ISBN10:048624895X",
		},
		want: true,
		err:  nil,
	}, {
		name: "isbn13",
		values: []string{
			"9780316129084",
			"ISBN13:9780316129084",
		},
		want: true,
		err:  nil,
	}, {
		name: "not isbn",
		values: []string{
			"foobar",
			"012345",
			"012345abc5",
			"012345abc5t12",
			"012345abcX",
			"abcdefghij",
			"Book Title ",
		},
		want: false,
		err:  nil,
	}, {
		name: "invalid isbn digits",
		values: []string{
			"123156789X",
			"1233445567891",
		},
		want: false,
		err:  ErrInvalidIsbn,
	}}

	is := is.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.values {
				got, err := IsbnCheck(v)
				if tt.err != nil {
					errors.Is(err, tt.err)
				} else {
					is.NoErr(err)
				}

				is.Equal(got, tt.want)
			}

		})
	}
}
