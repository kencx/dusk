package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/mock"
	"github.com/kencx/dusk/util"

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
	testHandler.db = &mock.Store{
		GetAuthorFn: func(id int64) (*dusk.Author, error) {
			return testAuthor1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/authors/1",
		params: map[string]string{"id": "1"},
		fn:     testHandler.GetAuthor,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["authors"]
	is.Equal(got.Name, testAuthor1.Name)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestGetAllAuthors(t *testing.T) {
	is := is.New(t)
	testHandler.db = &mock.Store{
		GetAllAuthorsFn: func() (dusk.Authors, error) {
			return testAuthors, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/authors/",
		fn:     testHandler.GetAllAuthors,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string][]*dusk.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["authors"]
	for i, v := range got {
		is.Equal(v.Name, testAuthors[i].Name)
	}
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestGetAllAuthorsNil(t *testing.T) {
	is := is.New(t)
	testHandler.db = &mock.Store{
		GetAllAuthorsFn: func() (dusk.Authors, error) {
			return nil, dusk.ErrNoRows
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/authors/",
		fn:     testHandler.GetAllAuthors,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	is.Equal(w.Code, http.StatusNoContent)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestAddAuthor(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testAuthor1)
	is.NoErr(err)

	testHandler.db = &mock.Store{
		CreateAuthorFn: func(a *dusk.Author) (*dusk.Author, error) {
			return testAuthor1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/authors/",
		data:   want,
		fn:     testHandler.AddAuthor,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["authors"]
	is.Equal(got.Name, testAuthor1.Name)
	is.Equal(w.Code, http.StatusCreated)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestAddAuthorFailValidation(t *testing.T) {
	is := is.New(t)
	failAuthor := &dusk.Author{Name: ""}
	want, err := util.ToJSON(failAuthor)
	is.NoErr(err)

	testHandler.db = &mock.Store{
		CreateAuthorFn: func(a *dusk.Author) (*dusk.Author, error) {
			return testAuthor1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/authors/",
		data:   want,
		fn:     testHandler.AddAuthor,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertValidationError(t, w, "name", "value is missing")
}

func TestUpdateAuthor(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testAuthor2)
	is.NoErr(err)

	testHandler.db = &mock.Store{
		UpdateAuthorFn: func(id int64, a *dusk.Author) (*dusk.Author, error) {
			return testAuthor2, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/authors/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testHandler.UpdateAuthor,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["authors"]
	is.Equal(got.Name, testAuthor2.Name)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestDeleteAuthor(t *testing.T) {
	is := is.New(t)
	testHandler.db = &mock.Store{
		DeleteAuthorFn: func(id int64) error {
			return nil
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/authors/1",
		params: map[string]string{"id": "1"},
		fn:     testHandler.DeleteAuthor,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)
	is.Equal(w.Code, http.StatusOK)
}
