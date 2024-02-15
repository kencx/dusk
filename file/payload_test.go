package file

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		mimetype string
		want     string
		err      error
	}{{
		name:     "success filename",
		filename: "foo.html",
		mimetype: "text/html",
		want:     ".html",
		err:      nil,
	}, {
		name:     "success known mimetype",
		filename: "foo",
		mimetype: "image/png",
		want:     ".png",
		err:      nil,
	}, {
		name:     "success custom mimetype",
		filename: "foo",
		mimetype: "application/epub+zip",
		want:     ".epub",
		err:      nil,
	}, {
		name:     "default mimetype",
		filename: "foo",
		mimetype: "foobar",
		want:     defaultMime,
		err:      nil,
	}}

	is := is.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extension(tt.filename, tt.mimetype)
			if tt.err == nil {
				is.NoErr(err)
			}

			if tt.err != nil {
				errors.Is(err, tt.err)
			}

			is.Equal(got, tt.want)
		})
	}
}
