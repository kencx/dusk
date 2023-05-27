package http

import (
	"bytes"
    "context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

var testServer = Server{
	InfoLog:  log.New(io.Discard, "", log.LstdFlags),
	ErrorLog: log.New(io.Discard, "", log.LstdFlags),
}

type testCase struct {
	url     string
	method  string
	headers map[string]string
	data    []byte
	params  map[string]string
	fn      func(http.ResponseWriter, *http.Request)
}

func testResponse(t *testing.T, tc *testCase) (*httptest.ResponseRecorder, error) {
	t.Helper()

	req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(tc.data))
	rw := httptest.NewRecorder()
    req = addTestParams(t, req, tc.params)

	http.HandlerFunc(tc.fn).ServeHTTP(rw, req)
	return rw, nil
}

func assertResponseError(t *testing.T, w *httptest.ResponseRecorder, status int, message string) {
	t.Helper()
	is := is.New(t)

	var env map[string]string
	err := json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["error"]
	is.Equal(w.Code, status)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json")
	is.Equal(got, message)
}

func assertValidationError(t *testing.T, w *httptest.ResponseRecorder, key, message string) {
	t.Helper()
	is := is.New(t)

	var env map[string]map[string]string
	err := json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["error"]
	is.Equal(w.Code, http.StatusUnprocessableEntity)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json")

	val, ok := got[key]
	if !ok {
		t.Errorf("validation error field %q not present", key)
	}
	is.Equal(val, message)
}

func addTestParams(t *testing.T, r *http.Request, params map[string]string) *http.Request {
    t.Helper()

    x := chi.NewRouteContext()
    routeParams := chi.RouteParams{}
    for k, v := range params {
        routeParams.Add(k, v)
    }
    x.URLParams = routeParams
    return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, x))
}
