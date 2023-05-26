package storage

import (
	"dusk"
	"dusk/util"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	testBook1 = &dusk.Book{
		ID:         1,
		Title:      "Book 1",
		ISBN:       "1",
		NumOfPages: 250,
		Rating:     5,
		State:      "read",
		Author:     []string{testAuthor1.Name},
		Tag:        []string{testTag1.Name},
	}
	testBook2 = &dusk.Book{
		ID:         2,
		Title:      "Book 2",
		ISBN:       "2",
		NumOfPages: 900,
		Rating:     4,
		State:      "unread",
		Author:     []string{testAuthor2.Name},
	}
	testBook3 = &dusk.Book{
		ID:     3,
		Title:  "Many Authors",
		ISBN:   "3",
		State:  "unread",
		Author: []string{testAuthor3.Name, testAuthor4.Name, testAuthor5.Name},
		Tag:    []string{testTag2.Name, testTag3.Name},
	}
	testBook4 = &dusk.Book{
		ID:     4,
		Title:  "Book 4",
		ISBN:   "4",
		State:  "unread",
		Author: []string{testAuthor5.Name},
	}
	allTestBooks = dusk.Books{testBook1, testBook2, testBook3, testBook4}
)

func TestGetBook(t *testing.T) {
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
		name: "not exists",
		id:   -1,
		want: nil,
		err:  dusk.ErrDoesNotExist,
	}}

	resetDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.GetBook(tt.id)
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}

			if !assertBooksEqual(got, tt.want) {
				t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(tt.want))
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
    got, err := ts.GetAllBooks()
    checkErr(t, err)

    want := allTestBooks

    if len(got) != len(want) {
        t.Errorf("got %d books, want %d books", len(got), len(want))
    }

    for i := 0; i < len(got); i++ {
        if !assertBooksEqual(got[i], want[i]) {
            t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
        }
    }
}

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name string
		want *dusk.Book
	}{{
		name: "book with minimal data",
		want: &dusk.Book{
			Title:  "1984",
			ISBN:   "1001",
			Author: []string{"George Orwell"},
		},
	}, {
		name: "book with all data",
		want: &dusk.Book{
			Title:      "World War Z",
			ISBN:       "1002",
			Author:     []string{"Max Brooks"},
			NumOfPages: 100,
			Rating:     10,
			State:      "read",
			Tag:        []string{"Zombies"},
		},
	}, {
		name: "book with two authors",
		want: &dusk.Book{
			Title:      "Pro Git",
			ISBN:       "1003",
			Author:     []string{"Scott Chacon", "Ben Straub"},
			NumOfPages: 100,
			Rating:     10,
			State:      "read",
		},
	}}

	defer resetDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := ts.CreateBook(tt.want)
			checkErr(t, err)

			got, err := ts.GetBook(b.ID)
			checkErr(t, err)

			if !assertBooksEqual(got, tt.want) {
				t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(tt.want))
			}

			assertAuthorsExist(t, tt.want)
			assertBookAuthorRelationship(t, tt.want)

			assertTagsExist(t, tt.want)
			assertBookTagRelationship(t, tt.want)
		})
	}
}

func TestCreateBookExistingISBN(t *testing.T) {
	_, err := ts.CreateBook(testBook2)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestCreateBookExistingAuthor(t *testing.T) {
    defer resetDB()
	want := &dusk.Book{
		Title:      "Morning Star",
		ISBN:       "1004",
		Author:     []string{"John Adams"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
	}
	got, err := ts.CreateBook(want)
	checkErr(t, err)

	assertAuthorsExist(t, want)
	assertBookAuthorRelationship(t, got)

    // TODO get authors related books
}

func TestCreateBookNewAndExistingAuthor(t *testing.T) {
	defer resetDB()
	want := &dusk.Book{
		Title:      "Tiamat's Wrath",
		ISBN:       "1005",
		Author:     []string{"John Adams", "Daniel Abrahams"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
	}
	_, err := ts.CreateBook(want)
	checkErr(t, err)

	assertAuthorsExist(t, want)
	assertBookAuthorRelationship(t, want)

    // TODO get both authors related books
}

func TestCreateBookExistingTag(t *testing.T) {
    defer resetDB()
	want := &dusk.Book{
		Title:      "Dune",
		ISBN:       "1008",
		Author:     []string{"Frank Herbert"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
		Tag:        []string{"Starred"},
	}
	_, err := ts.CreateBook(want)
	checkErr(t, err)

	assertTagsExist(t, want)
	assertBookTagRelationship(t, want)

    // TODO get tag's related books
}

func TestCreateBookNewAndExistingTag(t *testing.T) {
	defer resetDB()
	want := &dusk.Book{
		Title:      "Foundation",
		ISBN:       "1009",
		Author:     []string{"Isaac Asimov"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
		Tag:        []string{"New", "Starred"},
	}
	_, err := ts.CreateBook(want)
	checkErr(t, err)

	assertTagsExist(t, want)
	assertBookTagRelationship(t, want)

    // TODO get tag's related books
}

func TestUpdateBookNoAuthorChange(t *testing.T) {
	defer resetDB()

	want := testBook1
	want.NumOfPages = 999
	want.Rating = 1
	want.State = "unread"

	got, err := ts.UpdateBook(want.ID, want)
	checkErr(t, err)

	if !assertBooksEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
    // TODO get tag's related books
}

func TestUpdateBookAddNewAuthor(t *testing.T) {
	defer resetDB()

    want := modifyAuthors(testBook1, append(testBook1.Author, "Ty Franck"))
	got, err := ts.UpdateBook(want.ID, want)
	checkErr(t, err)

	if !assertBooksEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookAddExistingAuthor(t *testing.T) {
	defer resetDB()

    want := modifyAuthors(testBook1, append(testBook1.Author, testBook2.Author[0]))
	got, err := ts.UpdateBook(want.ID, want)
	checkErr(t, err)

	if !assertBooksEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookRemoveAuthor(t *testing.T) {
	defer resetDB()

    old := testBook3.Author
    want := modifyAuthors(testBook3, testBook3.Author[:len(testBook3.Author)-1])
	got, err := ts.UpdateBook(want.ID, want)
	checkErr(t, err)

	if !assertBooksEqual(got, want) {
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
	got, err := ts.UpdateBook(want.ID, want)
	checkErr(t, err)

	if !assertBooksEqual(got, want) {
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
	got, err := ts.UpdateBook(want.ID, want)
	checkErr(t, err)

	if !assertBooksEqual(got, want) {
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

func TestUpdateBookISBNConstraint(t *testing.T) {
	want := testBook1
	want.ISBN = testBook2.ISBN
	_, err := ts.UpdateBook(want.ID, want)
	if err == nil {
		t.Errorf("expected error: unique constraint ISBN")
	}
}


func TestDeleteBook(t *testing.T) {
	defer resetDB()
	err := ts.DeleteBook(testBook1.ID)
	checkErr(t, err)

	_, err = ts.GetBook(testBook1.ID)
	if err == nil {
		t.Errorf("expected error, book %d not deleted", testBook1.ID)
	}

	// check delete cascaded to book_author_link
	var dest []int
	stmt := `SELECT book FROM book_author_link WHERE book=$1`
	if err := ts.db.Select(&dest, stmt, testBook1.ID); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("deleting book %d did not cascade to book_author_link", testBook1.ID)
	}

	// check book's author deleted
	var destName []string
	stmt = `SELECT name FROM author WHERE id=$1`
	if err := ts.db.Select(&destName, stmt, testAuthor1.ID); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(destName) != 0 {
		t.Errorf("author %d not deleted", testAuthor1.ID)
	}
}

func TestDeleteBookEnsureAuthorRemainsForExistingBooks(t *testing.T) {
	defer resetDB()

	err := ts.DeleteBook(testBook3.ID)
	checkErr(t, err)

	// check delete cascaded to book_author_link
	var dest []int
	stmt := `SELECT book FROM book_author_link WHERE book=$1`
	if err := ts.db.Select(&dest, stmt, testBook3.ID); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("deleting book %d did not cascade to book_author_link", testBook1.ID)
	}

	// check author still exists in authors table
	query, args, err := sqlx.In("SELECT id FROM author WHERE name IN (?);", testBook3.Author)
	checkErr(t, err)

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
	checkErr(t, err)

	var bookID []int
	query = ts.db.Rebind(query)
	if err := ts.db.Select(&bookID, query, args...); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(bookID) != 1 {
		t.Errorf("number of linked books incorrect")
	}

	got, err := ts.GetBook(testBook4.ID)
	checkErr(t, err)
	if !assertBooksEqual(got, testBook4) {
		t.Errorf("got %v, want %v", got, testBook4)
	}
}

func TestDeleteBookNotExists(t *testing.T) {
	err := ts.DeleteBook(-1)
	if err == nil {
		t.Errorf("expected error: book not exists")
	}
}

func modifyAuthors(testBook *dusk.Book, new []string) *dusk.Book {
    newBook := *testBook
    newBook.Author = new
    return &newBook
}

func assertAuthorsExist(t *testing.T, want *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
        authors, err := getAuthorsFromBook(tx, want.ID)
        checkErr(t, err)

        util.Sort(want.Author)
        if !reflect.DeepEqual(authors, want.Author) {
            t.Errorf("got %v, want %v", authors, want.Author)
        }
        return nil, nil
    })
}

func assertTagsExist(t *testing.T, want *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
        tags, err := getTagsFromBook(tx, want.ID)
        checkErr(t, err)

        util.Sort(want.Tag)
        if !reflect.DeepEqual(tags, want.Tag) {
            t.Errorf("got %v, want %v", tag, want.Tag)
        }
        return nil, nil
    })
}

func assertBookAuthorRelationship(t *testing.T, book *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
		// get book's related authors
		authors, err := getModelsFromBook(tx, book.ID, author)
		checkErr(t, err)

		if len(authors) != len(book.Author) {
			t.Errorf("book has wrong number of authors in books_authors table")
		}

		util.Sort(book.Author)

		// author must exist in book's related authors
		if !reflect.DeepEqual(authors, book.Author) {
			t.Errorf("got %v, want %v", authors, book.Author)
		}
		return nil, nil
	})
}

func assertBookTagRelationship(t *testing.T, book *dusk.Book) {
	t.Helper()
	Tx(ts.db, func(tx *sqlx.Tx) (any, error) {
		// get book's related tags
		tags, err := getModelsFromBook(tx, book.ID, tag)
		checkErr(t, err)

		if len(tags) != len(book.Tag) {
			t.Errorf("book has wrong number of tags in books_tags table")
		}

		util.Sort(book.Author)

		// tag must exist in book's related tags
		if !reflect.DeepEqual(tags, book.Tag) {
			t.Errorf("got %v, want %v", tags, book.Tag)
		}
		return nil, nil
	})
}

func assertBooksEqual(a, b *dusk.Book) bool {
	if (a == nil) && (b == nil) {
		return true
	}
	if (a != nil) && (b != nil) {
		authorEqual := reflect.DeepEqual(a.Author, b.Author)
		tagEqual := reflect.DeepEqual(a.Tag, b.Tag)
		return (a.Title == b.Title &&
			a.ISBN == b.ISBN &&
			a.NumOfPages == b.NumOfPages &&
			a.State == b.State &&
			a.Rating == b.Rating &&
			authorEqual &&
			tagEqual)
	}
	return a == b
}
