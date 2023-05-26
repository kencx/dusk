package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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
	if tc.params != nil {
		req = mux.SetURLVars(req, tc.params)
	}

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
