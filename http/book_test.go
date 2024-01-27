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
	testBook1 = &dusk.Book{
		Title:  "Book 1",
		Author: []string{"John Adams"},
		ISBN:   "100",
	}
	testBook2 = &dusk.Book{
		Title:  "Book 2",
		Author: []string{"Alice Brown"},
		ISBN:   "101",
	}
	testBook3 = &dusk.Book{
		Title:  "Book 3",
		Author: []string{"Billy Foo", "Carl Baz"},
		ISBN:   "102",
	}
	testBooks = []*dusk.Book{testBook1, testBook2, testBook3}
)

func TestGetBook(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		GetBookFn: func(id int64) (*dusk.Book, error) {
			return testBook1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/1",
		params: map[string]string{"id": "1"},
		fn:     testServer.GetBook,
	}
	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["books"]
	is.Equal(got.Title, testBook1.Title)
	is.Equal(got.Author[0], testBook1.Author[0])
	is.Equal(got.ISBN, testBook1.ISBN)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestGetBookNil(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		GetBookFn: func(id int64) (*dusk.Book, error) {
			return nil, dusk.ErrDoesNotExist
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/1",
		params: map[string]string{"id": "1"},
		fn:     testServer.GetBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertResponseError(t, w, http.StatusNotFound, "the item does not exist")
}

func TestGetAllBooks(t *testing.T) {
	is := is.New(t)
	t.Run("success", func(t *testing.T) {
		testServer.db = &mock.Store{
			GetAllBooksFn: func() (dusk.Books, error) {
				return testBooks, nil
			},
		}

		tc := &testCase{
			method: http.MethodGet,
			url:    "/api/books/",
			fn:     testServer.GetAllBooks,
		}

		w, err := testResponse(t, tc)
		is.NoErr(err)

		var env map[string][]*dusk.Book
		err = json.NewDecoder(w.Body).Decode(&env)
		is.NoErr(err)

		got := env["books"]
		for i, v := range got {
			is.Equal(v.Title, testBooks[i].Title)
			is.Equal(v.Author[0], testBooks[i].Author[0])
			is.Equal(v.ISBN, testBooks[i].ISBN)
		}
		is.Equal(w.Code, http.StatusOK)
		is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
	})

	t.Run("no content", func(t *testing.T) {
		testServer.db = &mock.Store{
			GetAllBooksFn: func() (dusk.Books, error) {
				return nil, dusk.ErrNoRows
			},
		}

		tc := &testCase{
			method: http.MethodGet,
			url:    "/api/books/",
			fn:     testServer.GetAllBooks,
		}

		w, err := testResponse(t, tc)
		is.NoErr(err)
		is.Equal(w.Code, http.StatusNoContent)
	})
}

func TestAddBook(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testBook1)
	is.NoErr(err)

	testServer.db = &mock.Store{
		CreateBookFn: func(b *dusk.Book) (*dusk.Book, error) {
			return testBook1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/books",
		data:   want,
		fn:     testServer.AddBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["books"]
	is.Equal(got.Title, testBook1.Title)
	is.Equal(got.Author[0], testBook1.Author[0])
	is.Equal(got.ISBN, testBook1.ISBN)
	is.Equal(w.Code, http.StatusCreated)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestAddBookFailValidation(t *testing.T) {
	is := is.New(t)
	failBook := &dusk.Book{
		Title:  "",
		Author: []string{"John Doe"},
		ISBN:   "12345",
	}
	want, err := util.ToJSON(failBook)
	is.NoErr(err)

	testServer.db = &mock.Store{
		CreateBookFn: func(b *dusk.Book) (*dusk.Book, error) {
			return failBook, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/books",
		data:   want,
		fn:     testServer.AddBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertValidationError(t, w, "title", "value is missing")
}

func TestUpdateBook(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testBook2)
	is.NoErr(err)

	testServer.db = &mock.Store{
		UpdateBookFn: func(id int64, b *dusk.Book) (*dusk.Book, error) {
			return testBook2, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/books/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testServer.UpdateBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)

	var env map[string]*dusk.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	is.NoErr(err)

	got := env["books"]
	is.Equal(got.Title, testBook2.Title)
	is.Equal(got.Author[0], testBook2.Author[0])
	is.Equal(got.ISBN, testBook2.ISBN)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Result().Header.Get("Content-Type"), "application/json")
}

func TestUpdateBookNil(t *testing.T) {
	is := is.New(t)
	want, err := util.ToJSON(testBook2)
	is.NoErr(err)

	testServer.db = &mock.Store{
		UpdateBookFn: func(id int64, b *dusk.Book) (*dusk.Book, error) {
			return nil, dusk.ErrDoesNotExist
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/books/10",
		data:   want,
		params: map[string]string{"id": "10"},
		fn:     testServer.UpdateBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertResponseError(t, w, http.StatusNotFound, "the item does not exist")
}

func TestUpdateBookFailValidation(t *testing.T) {
	is := is.New(t)
	failBook := &dusk.Book{
		Title:  "",
		Author: []string{"John Doe"},
		ISBN:   "12345",
	}
	want, err := util.ToJSON(failBook)
	is.NoErr(err)

	testServer.db = &mock.Store{
		UpdateBookFn: func(id int64, b *dusk.Book) (*dusk.Book, error) {
			return failBook, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/books/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testServer.UpdateBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertValidationError(t, w, "title", "value is missing")
}

func TestDeleteBook(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		DeleteBookFn: func(id int64) error {
			return nil
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/books/1",
		params: map[string]string{"id": "1"},
		fn:     testServer.DeleteBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)
	is.Equal(w.Code, http.StatusOK)
}

func TestDeleteBookNil(t *testing.T) {
	is := is.New(t)
	testServer.db = &mock.Store{
		DeleteBookFn: func(id int64) error {
			return dusk.ErrDoesNotExist
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/books/10",
		params: map[string]string{"id": "10"},
		fn:     testServer.DeleteBook,
	}

	w, err := testResponse(t, tc)
	is.NoErr(err)
	assertResponseError(t, w, http.StatusNotFound, "the item does not exist")
}
