package http

import (
	"dusk"
	"dusk/mock"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/matryer/is"
)

var (
	testAuthor1 = &dusk.Author{
		Name: "Author 1",
	}
	testAuthor2 = &dusk.Author{
		Name: "Author 2",
	}
	testAuthors = []*dusk.Author{testAuthor1, testAuthor2}
)

func TestGetAuthor(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		GetAuthorFn: func(id int64) (*dusk.Author, error) {
			return testAuthor1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/authors/1",
		params: map[string]string{"id": "1"},
		fn:     testServer.GetAuthor,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["authors"]
	is.Equal(got.Name, testAuthor1.Name)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json")
}
