package storage

import (
	"testing"
	"time"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/util"

	"github.com/jmoiron/sqlx"
	"github.com/matryer/is"
)

func TestGetBook(t *testing.T) {
	resetDB()
	tests := []struct {
		name string
		id   int64
		want *dusk.Book
		err  error
	}{{
		name: "book 2",
		id:   2,
		want: testBook2,
		err:  nil,
	}, {
		name: "book 3",
		id:   3,
		want: testBook3,
		err:  nil,
	}, {
		name: "book not exists",
		id:   -1,
		want: nil,
		err:  dusk.ErrDoesNotExist,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.GetBook(tt.id)
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}

			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(tt.want))
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetAllBooks()
	is.NoErr(err)

	want := allTestBooks

	if len(got) != len(want) {
		t.Errorf("got %d books, want %d books", len(got), len(want))
	}

	for i := 0; i < len(got); i++ {
		if !got[i].Equal(want[i]) {
			t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
		}
	}
}

// func TestGetAllByTitle(t *testing.T) {
// 	is := is.New(t)
//
// 	t.Run("successful full query", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Title = testBook1.Title
//
// 		got, _, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := []*dusk.Book{testBook1}
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
//
// 	t.Run("successful FTS query", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Title = strings.Split(testBook1.Title, " ")[0]
//
// 		got, _, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := []*dusk.Book{testBook1}
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
//
// 	t.Run("no match query", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Title = "foo"
// 		got, _, err := ts.GetAllBooks(filters)
//
// 		if err == nil {
// 			t.Fatalf("expected err: ErrDoesNotExist")
// 		}
//
// 		if err != dusk.ErrNoRows {
// 			t.Fatalf("unexpected err: %v", err)
// 		}
//
// 		if len(got) != 0 {
// 			t.Fatalf("got %v books, want %v books", got, 0)
// 		}
// 	})
// }
//
// func TestGetAllByAuthor(t *testing.T) {
// 	is := is.New(t)
//
// 	t.Run("successful query", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Author = testAuthor1.Name
//
// 		got, _, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := []*dusk.Book{testBook1}
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
//
// 	t.Run("successful FTS query", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Author = strings.Split(testAuthor1.Name, " ")[1]
//
// 		got, _, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := []*dusk.Book{testBook1}
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
//
// 	t.Run("no match query", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Author = "foo"
//
// 		got, _, err := ts.GetAllBooks(filters)
// 		if err == nil {
// 			t.Fatalf("expected err: ErrDoesNotExist")
// 		}
//
// 		if err != dusk.ErrNoRows {
// 			t.Fatalf("unexpected err: %v", err)
// 		}
//
// 		if len(got) != 0 {
// 			t.Fatalf("got %v books, want %v books", got, 0)
// 		}
// 	})
// }
//
// func TestGetAllSortBy(t *testing.T) {
// 	t.Run("sort by title DESC", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.Sort = "-title"
//
// 		got, _, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := allTestBooks
// 		sort.SliceStable(want, func(i, j int) bool {
// 			return want[j].Title < want[i].Title
// 		})
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
// }
//
// func TestGetAllPagination(t *testing.T) {
// 	is := is.New(t)
//
// 	t.Run("success", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.AfterId = 2
// 		filters.PageSize = 2
//
// 		got, count, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := []*dusk.Book{testBook3, testBook4}
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		if count != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
//
// 	t.Run("last entry only", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.AfterId = len(allTestBooks) - 1
//
// 		got, count, err := ts.GetAllBooks(filters)
// 		is.NoErr(err)
//
// 		want := []*dusk.Book{testBook4}
//
// 		if len(got) != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		if count != len(want) {
// 			t.Fatalf("got %d books, want %d books", len(got), len(want))
// 		}
//
// 		for i := 0; i < len(got); i++ {
// 			if !got[i].Equal(want[i]) {
// 				t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 			}
// 		}
// 	})
//
// 	t.Run("no more results", func(t *testing.T) {
// 		filters := testFilters()
// 		filters.AfterId = len(allTestBooks)
//
// 		got, count, err := ts.GetAllBooks(filters)
// 		if err == nil {
// 			t.Fatalf("expected err: ErrDoesNotExist")
// 		}
//
// 		if err != dusk.ErrNoRows {
// 			t.Fatalf("unexpected err: %v", err)
// 		}
//
// 		if len(got) != 0 {
// 			t.Fatalf("got %v books, want %v books", got, 0)
// 		}
//
// 		if count != 0 {
// 			t.Fatalf("got %v books, want %v books", got, 0)
// 		}
// 	})
// }
//
// func TestGetAllMultipleFilters(t *testing.T) {
// 	is := is.New(t)
//
// 	filters := testFilters()
// 	filters.AfterId = 1
// 	filters.Title = strings.Split(testBook3.Title, " ")[0]
// 	filters.Author = strings.Split(testAuthor3.Name, " ")[0]
//
// 	got, count, err := ts.GetAllBooks(filters)
// 	is.NoErr(err)
//
// 	want := []*dusk.Book{testBook3}
//
// 	if len(got) != len(want) {
// 		t.Fatalf("got %d books, want %d books", len(got), len(want))
// 	}
//
// 	if count != len(want) {
// 		t.Fatalf("got %d books, want %d books", len(got), len(want))
// 	}
//
// 	for i := 0; i < len(got); i++ {
// 		if !got[i].Equal(want[i]) {
// 			t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
// 		}
// 	}
// }

func TestCreateBook(t *testing.T) {
	defer resetDB()

	tests := []struct {
		name string
		want *dusk.Book
	}{{
		name: "simple book",
		want: &dusk.Book{
			Title:  "1984",
			Author: []string{"George Orwell"},
		},
	}, {
		name: "full book",
		want: &dusk.Book{
			Title:         "Book 6",
			Subtitle:      null.StringFrom("subtitle 1"),
			Isbn10:        []string{"1451673310"},
			Author:        []string{"author 6"},
			NumOfPages:    100,
			Rating:        10,
			Progress:      80,
			Tag:           []string{"tag 4"},
			Publisher:     null.StringFrom("publisher 1"),
			DatePublished: null.TimeFrom(time.Now()),
			Series:        null.StringFrom("series 2"),
			Description:   null.StringFrom("lorem ipsum"),
			DateStarted:   null.TimeFrom(time.Now()),
			DateCompleted: null.TimeFrom(time.Now()),
		},
	}, {
		name: "book with two authors",
		want: &dusk.Book{
			Title:   "Book 7",
			Author:  []string{"author 7", "author 8"},
			Tag:     []string{"tag 5", "tag 6"},
			Formats: []string{"format 2", "format 3"},
		},
	}}

	is := is.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := ts.CreateBook(tt.want)
			is.NoErr(err)

			got, err := ts.GetBook(b.Id)
			is.NoErr(err)

			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(tt.want))
			}

			assertAuthorsExist(t, got)
			assertTagsExist(t, got)
		})
	}
}

func TestCreateBookExistingIsbn10(t *testing.T) {
	_, err := ts.CreateBook(testBook2)
	if err == nil {
		t.Errorf("expected error: unique constraint Isbn10")
	}
}

func TestCreateBookExistingIsbn13(t *testing.T) {
	_, err := ts.CreateBook(testBook4)
	if err == nil {
		t.Errorf("expected error: unique constraint Isbn13")
	}
}

func TestCreateBookExistingAuthor(t *testing.T) {
	defer resetDB()

	want := &dusk.Book{
		Title:  "Book 8",
		Author: []string{testAuthor1.Name},
	}
	is := is.New(t)
	got, err := ts.CreateBook(want)
	is.NoErr(err)

	assertAuthorsExist(t, got)

	relatedBooks, err := ts.GetAllBooksFromAuthor(testAuthor1.Id)
	is.NoErr(err)
	if len(relatedBooks) != 2 {
		t.Errorf("got %d books, want %d books", len(relatedBooks), 2)
	}
}

func TestCreateBookNewAndExistingAuthor(t *testing.T) {
	defer resetDB()
	want := &dusk.Book{
		Title:  "Book 9",
		Author: []string{testAuthor1.Name, "author 9"},
	}
	is := is.New(t)
	got, err := ts.CreateBook(want)
	is.NoErr(err)

	assertAuthorsExist(t, got)

	relatedBooks, err := ts.GetAllBooksFromAuthor(testAuthor1.Id)
	is.NoErr(err)
	if len(relatedBooks) != 2 {
		t.Errorf("got %d books, want %d books", len(relatedBooks), 2)
	}
}

func TestCreateBookExistingTag(t *testing.T) {
	defer resetDB()
	want := &dusk.Book{
		Title:  "Book 10",
		Author: []string{"author 10"},
		Tag:    []string{testTag1.Name},
	}
	is := is.New(t)
	got, err := ts.CreateBook(want)
	is.NoErr(err)

	assertTagsExist(t, got)

	relatedBooks, err := ts.GetAllBooksFromTag(testTag1.Id)
	is.NoErr(err)
	if len(relatedBooks) != 2 {
		t.Errorf("got %d books, want %d books", len(relatedBooks), 2)
	}
}

func TestCreateBookNewAndExistingTag(t *testing.T) {
	defer resetDB()
	want := &dusk.Book{
		Title:  "Book 11",
		Author: []string{"author 11"},
		Tag:    []string{"tag 4", testTag1.Name},
	}
	is := is.New(t)
	got, err := ts.CreateBook(want)
	is.NoErr(err)

	assertTagsExist(t, got)
}

func TestCreateBookExistingSeries(t *testing.T) {
	defer resetDB()

	want := &dusk.Book{
		Title:  "Book 12",
		Author: []string{"author 12"},
		Series: null.StringFrom("series 1"),
	}

	is := is.New(t)
	got, err := ts.CreateBook(want)
	is.NoErr(err)

	assertSeriesExist(t, got)
}

func TestUpdateBookNoAuthorChange(t *testing.T) {
	defer resetDB()

	want := *testBook1
	want.NumOfPages = 999
	want.Rating = 1

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, &want)
	is.NoErr(err)

	if !got.Equal(&want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestUpdateBookAddNewAuthor(t *testing.T) {
	defer resetDB()

	want := modifyAuthors(testBook1, append(testBook1.Author, "Ty Franck"))
	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, want)
	is.NoErr(err)

	if !got.Equal(want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)
}

func TestUpdateBookAddExistingAuthor(t *testing.T) {
	defer resetDB()

	want := modifyAuthors(testBook1, append(testBook1.Author, testBook2.Author[0]))
	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, want)
	is.NoErr(err)

	if !got.Equal(want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)
}

func TestUpdateBookRemoveAuthor(t *testing.T) {
	defer resetDB()

	old := testBook3.Author
	want := modifyAuthors(testBook3, testBook3.Author[:len(testBook3.Author)-1])
	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, want)
	is.NoErr(err)

	if !got.Equal(want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check removed author still exists
	var dest []string
	last := old[len(old)-1]
	stmt := `SELECT name FROM author WHERE name=$1`
	if err := ts.db.Select(&dest, stmt, last); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	if len(dest) == 0 {
		t.Errorf("author %s missing", last)
	}

	// check relationship with previous author dropped
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookRemoveAuthorCompletely(t *testing.T) {
	defer resetDB()

	old := testBook2.Author[0]
	want := modifyAuthors(testBook2, testBook1.Author)
	want.Rating = 2

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, want)
	is.NoErr(err)

	if !got.Equal(want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check author removed permanently
	var dest []string
	stmt := `SELECT name FROM author WHERE name=$1`
	if err := ts.db.Select(&dest, stmt, old); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	if len(dest) != 0 {
		t.Errorf("author %s still exists", old)
	}

	// check relationship with previous author dropped
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookRenameAuthor(t *testing.T) {
	defer resetDB()

	old := testBook4.Author[0]
	want := modifyAuthors(testBook4, []string{"Daniel Foo"})
	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, want)
	is.NoErr(err)

	if !got.Equal(want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)

	// check author still exists
	var dest []string
	stmt := `SELECT name FROM author WHERE name=$1`
	if err := ts.db.Select(&dest, stmt, old); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	if len(dest) == 0 {
		t.Errorf("author %s does not exist", old)
	}

	// relationship with previous author dropped
	// new relationship formed
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookNotExists(t *testing.T) {
	b := &dusk.Book{}
	_, err := ts.UpdateBook(-1, b)
	if err == nil {
		t.Errorf("expected error: no books updated")
	}
}

func TestUpdateBookIsbn10Constraint(t *testing.T) {
	want := *testBook1
	want.Isbn10 = testBook2.Isbn10
	_, err := ts.UpdateBook(want.Id, &want)
	if err == nil {
		t.Errorf("expected error: unique constraint Isbn10")
	}
}

func TestUpdateBookIsbn13Constraint(t *testing.T) {
	want := *testBook2
	want.Isbn13 = testBook4.Isbn13
	_, err := ts.UpdateBook(want.Id, &want)
	if err == nil {
		t.Errorf("expected error: unique constraint Isbn13")
	}
}

func TestUpdateBookTags(t *testing.T) {
	defer resetDB()

	old := testBook1.Tag[0]
	want := *testBook1
	want.Tag = []string{"tag 5"}

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, &want)
	is.NoErr(err)

	if !got.Equal(&want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertTagsExist(t, &want)

	// check old tag still exists
	var dest []string
	stmt := `SELECT name FROM tag WHERE name=$1`
	if err := ts.db.Select(&dest, stmt, old); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	if len(dest) == 0 {
		t.Errorf("tag %s does not exist", old)
	}

	// relationship with previous tag dropped
	// new relationship formed
	assertBookTagRelationship(t, &want)
}

func TestUpdateBookIsbn10(t *testing.T) {
	defer resetDB()

	old := testBook2.Isbn10[0]
	want := *testBook2
	want.Isbn10 = []string{"0547928211"}

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, &want)
	is.NoErr(err)

	if !got.Equal(&want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check old isbn does not exist
	var dest []string
	stmt := `SELECT isbn FROM isbn10 WHERE bookId=$1`
	if err := ts.db.Select(&dest, stmt, want.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	for _, i := range dest {
		if i == old {
			t.Errorf("isbn10 %s still exists", old)
		}
	}
}

func TestUpdateBookIsbn13(t *testing.T) {
	defer resetDB()

	old := testBook4.Isbn13[0]
	want := *testBook4
	want.Isbn13 = []string{"9780547928210"}

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, &want)
	is.NoErr(err)

	if !got.Equal(&want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check old isbn does not exist
	var dest []string
	stmt := `SELECT isbn FROM isbn13 WHERE bookId=$1`
	if err := ts.db.Select(&dest, stmt, want.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	for _, i := range dest {
		if i == old {
			t.Errorf("isbn13 %s still exists", old)
		}
	}
}

func TestUpdateBookSeries(t *testing.T) {
	defer resetDB()

	old := testBook2.Series.ValueOrZero()
	want := *testBook2
	want.Series = null.StringFrom("series 3")

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, &want)
	is.NoErr(err)

	if !got.Equal(&want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check old series does not exist
	var dest []string
	stmt := `SELECT name FROM series WHERE bookId=$1`
	if err := ts.db.Select(&dest, stmt, want.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	for _, i := range dest {
		if i == old {
			t.Errorf("series %s still exists", old)
		}
	}
}

func TestUpdateBookFormat(t *testing.T) {
	defer resetDB()

	old := testBook2.Formats[0]
	want := *testBook2
	want.Formats = []string{"format 4"}

	is := is.New(t)
	got, err := ts.UpdateBook(want.Id, &want)
	is.NoErr(err)

	if !got.Equal(&want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check old format does not exist
	var dest []string
	stmt := `SELECT filepath FROM format WHERE bookId=$1`
	if err := ts.db.Select(&dest, stmt, want.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	for _, i := range dest {
		if i == old {
			t.Errorf("format %s still exists", old)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	err := ts.DeleteBook(testBook1.Id)
	is.NoErr(err)

	_, err = ts.GetBook(testBook1.Id)
	if err == nil {
		t.Errorf("expected error, book %d not deleted", testBook1.Id)
	}

	// check delete cascaded to book_author_link
	var dest []int
	stmt := `SELECT book FROM book_author_link WHERE book=$1`
	if err := ts.db.Select(&dest, stmt, testBook1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("deleting book %d did not cascade to book_author_link", testBook1.Id)
	}

	// check book's author deleted
	var destName []string
	stmt = `SELECT name FROM author WHERE id=$1`
	if err := ts.db.Select(&destName, stmt, testAuthor1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) != 0 {
		t.Errorf("author %d not deleted", testAuthor1.Id)
	}

	// check delete cascaded to book_tag_link
	stmt = `SELECT book FROM book_tag_link WHERE book=$1`
	if err := ts.db.Select(&dest, stmt, testBook1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("deleting book %d did not cascade to book_tag_link", testBook1.Id)
	}

	// check book's tags not deleted
	stmt = `SELECT name FROM tag WHERE id=$1`
	if err := ts.db.Select(&destName, stmt, testTag1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) == 0 {
		t.Errorf("tag %d deleted", testTag1.Id)
	}

	// check isbn10 deleted
	stmt = `SELECT isbn FROM isbn10 WHERE bookId=$1`
	if err := ts.db.Select(&destName, stmt, testBook1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) != 0 {
		t.Errorf("isbn10 %s not deleted", destName[0])
	}

	// check isbn13 deleted
	stmt = `SELECT isbn FROM isbn13 WHERE bookId=$1`
	if err := ts.db.Select(&destName, stmt, testBook1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) != 0 {
		t.Errorf("isbn13 %s not deleted", destName[0])
	}

	// check series deleted
	stmt = `SELECT name FROM series WHERE bookId=$1`
	if err := ts.db.Select(&destName, stmt, testBook1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) != 0 {
		t.Errorf("series %s not deleted", destName[0])
	}

	// check format deleted
	stmt = `SELECT filepath FROM format WHERE bookId=$1`
	if err := ts.db.Select(&destName, stmt, testBook1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) != 0 {
		t.Errorf("format %s not deleted", destName[0])
	}
}

func TestDeleteBookEnsureAuthorRemainsForExistingBooks(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	err := ts.DeleteBook(testBook3.Id)
	is.NoErr(err)

	// check delete cascaded to book_author_link
	var dest []int
	stmt := `SELECT book FROM book_author_link WHERE book=$1`
	if err := ts.db.Select(&dest, stmt, testBook3.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("deleting book %d did not cascade to book_author_link", testBook1.Id)
	}

	// check author still exists in authors table
	query, args, err := sqlx.In("SELECT id FROM author WHERE name IN (?);", testBook3.Author)
	is.NoErr(err)

	var count []int
	query = ts.db.Rebind(query)
	if err = ts.db.Select(&count, query, args...); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(count) != 1 {
		t.Errorf("got %d author, want 1 author", len(count))
	}

	// check if author still linked to their other books in book_author_link
	stmt = `SELECT ba.book
		FROM book_author_link ba
            JOIN author a ON a.id=ba.author
		WHERE a.name IN (?);`
	query, args, err = sqlx.In(stmt, testBook3.Author)
	is.NoErr(err)

	var bookId []int
	query = ts.db.Rebind(query)
	if err := ts.db.Select(&bookId, query, args...); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(bookId) != 1 {
		t.Errorf("number of linked books incorrect")
	}

	got, err := ts.GetBook(testBook4.Id)
	is.NoErr(err)

	if !got.Equal(testBook4) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(testBook4))
	}
}

func TestDeleteBookNotExists(t *testing.T) {
	err := ts.DeleteBook(-1)
	if err == nil {
		t.Errorf("expected error: book not exists")
	}
}

func modifyAuthors(old *dusk.Book, newAuthors []string) *dusk.Book {
	newBook := *old
	newBook.Author = newAuthors
	return &newBook
}

func assertAuthorsExist(t *testing.T, b *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
		is := is.New(t)
		authors, err := getAuthorsFromBook(tx, b.Id)
		is.NoErr(err)

		util.Sort(b.Author)
		is.Equal(authors, b.Author)
		return nil, nil
	})
}

func assertTagsExist(t *testing.T, b *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
		is := is.New(t)
		tags, err := getTagsFromBook(tx, b.Id)
		is.NoErr(err)

		util.Sort(b.Tag)
		is.Equal(tags, b.Tag)
		return nil, nil
	})
}

func assertBookAuthorRelationship(t *testing.T, b *dusk.Book) {
	assertAuthorsExist(t, b)
}

func assertBookTagRelationship(t *testing.T, b *dusk.Book) {
	assertTagsExist(t, b)
}

func assertSeriesExist(t *testing.T, want *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
		is := is.New(t)
		series, err := getSeriesFromBook(tx, want.Id)
		is.NoErr(err)

		is.Equal(series.Name, want.Series.ValueOrZero())
		return nil, nil
	})
}
