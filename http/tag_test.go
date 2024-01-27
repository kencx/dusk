package http

import (
	"dusk"
	"dusk/mock"
	"dusk/util"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/matryer/is"
)

var (
	testTag1 = &dusk.Tag{
		Name: "Tag 1",
	}
	testTag2 = &dusk.Tag{
		Name: "Tag 2",
	}
	testTags = []*dusk.Tag{testTag1, testTag2}
)

func TestGetTag(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		GetTagFn: func(id int64) (*dusk.Tag, error) {
			return testTag1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/tags/1",
		params: map[string]string{"id": "1"},
		fn:     testServer.GetTag,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Tag
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["tags"]
	is.Equal(got.Name, testTag1.Name)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestGetAllTags(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		GetAllTagsFn: func() (dusk.Tags, error) {
			return testTags, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/tags/",
		fn:     testServer.GetAllTags,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string][]*dusk.Tag
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["tags"]
	for i, v := range got {
		is.Equal(v.Name, testTags[i].Name)
	}
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestGetAllTagsNil(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		GetAllTagsFn: func() (dusk.Tags, error) {
			return nil, dusk.ErrNoRows
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/tags/",
		fn:     testServer.GetAllTags,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	is.Equal(w.Code, http.StatusNoContent)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestAddTag(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testTag1)
	is.NoErr(err)

	testServer.db = &mock.Store{
		CreateTagFn: func(a *dusk.Tag) (*dusk.Tag, error) {
			return testTag1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/tags/",
		data:   want,
		fn:     testServer.AddTag,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Tag
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["tags"]
	is.Equal(got.Name, testTag1.Name)
	is.Equal(w.Code, http.StatusCreated)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestAddTagFailValidation(t *testing.T) {
	is := is.New(t)
	failTag := &dusk.Tag{Name: ""}
	want, err := util.ToJSON(failTag)
	is.NoErr(err)

	testServer.db = &mock.Store{
		CreateTagFn: func(a *dusk.Tag) (*dusk.Tag, error) {
			return testTag1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/tags/",
		data:   want,
		fn:     testServer.AddTag,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertValidationError(t, w, "name", "value is missing")
}

func TestUpdateTag(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testTag2)
	is.NoErr(err)

	testServer.db = &mock.Store{
		UpdateTagFn: func(id int64, a *dusk.Tag) (*dusk.Tag, error) {
			return testTag2, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/tags/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testServer.UpdateTag,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Tag
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["tags"]
	is.Equal(got.Name, testTag2.Name)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestDeleteTag(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		DeleteTagFn: func(id int64) error {
			return nil
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/tags/1",
		params: map[string]string{"id": "1"},
		fn:     testServer.DeleteTag,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)
	is.Equal(w.Code, http.StatusOK)
}
