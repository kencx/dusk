package storage

import (
	"testing"

	"github.com/kencx/dusk"

	"github.com/matryer/is"
)

func TestGetAuthor(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetAuthor(testAuthor1.Id)
	is.NoErr(err)

	want := testAuthor1
	is.Equal(got, want)
}

func TestGetAuthorNotExists(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetAuthor(-1)
	if err == nil {
		t.Errorf("expected error: ErrDoesNotExist")
	}

	if err != dusk.ErrDoesNotExist {
		t.Errorf("unexpected error: %v", err)
	}
	is.Equal(got, nil)
}

// func TestGetAuthorWithName(t *testing.T) {
// 	got, err := ts.Authors.GetByName(testAuthor2.Name)
// 	checkErr(t, err)
//
// 	want := testAuthor2
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
// 	}
// }

func TestGetAllAuthors(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetAllAuthors()
	is.NoErr(err)

	want := allTestAuthors
	is.Equal(got, want)
}

func TestGetAllAuthorEmpty(t *testing.T) {
	defer resetDB()

	// delete all data
	if err := ts.MigrateUp(resetSchemaPath); err != nil {
		t.Errorf("failed to reset database")
	}

	got, err := ts.GetAllAuthors()
	if err == nil {
		t.Errorf("expected error: ErrNoRows")
	}

	if err != dusk.ErrNoRows {
		t.Errorf("unexpected error: %v", err)
	}

	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestGetAllBooksFromAuthor(t *testing.T) {
	is := is.New(t)

	got, err := ts.GetAllBooksFromAuthor(testAuthor5.Id)
	is.NoErr(err)

	want := dusk.Books{testBook3, testBook4}
	is.True(len(got) == len(want))
	for i := range want {
		is.Equal(got[i].Title, want[i].Title)
	}
}

func TestCreateAuthor(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	want := &dusk.Author{Name: "FooBar"}

	got, err := ts.CreateAuthor(want)
	is.NoErr(err)
	is.Equal(got, want)
}

func TestCreateAuthorDuplicates(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	want := testAuthor3

	got, err := ts.CreateAuthor(want)
	is.NoErr(err)
	is.Equal(got, want)

	// check for number of entries in authors
	var dest []string
	stmt := `SELECT name FROM author WHERE name=$1`
	if err := ts.db.Select(&dest, stmt, want.Name); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(dest) != 1 {
		t.Error("more than one author inserted")
	}
}

func TestUpdateAuthor(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	want := *testAuthor1
	want.Name = "Sherlock Holmes"

	got, err := ts.UpdateAuthor(want.Id, &want)
	is.NoErr(err)
	is.Equal(got.Name, want.Name)
}

func TestUpdateAuthorExisting(t *testing.T) {
	want := *testAuthor1
	want.Name = testAuthor2.Name

	_, err := ts.UpdateAuthor(want.Id, &want)
	if err == nil {
		t.Errorf("expected error: unique constraint Name")
	}
}

func TestDeleteAuthor(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	// delete book first to circumvent foreign key constraint
	stmt := `DELETE from book WHERE id=$1;`
	_, err := ts.db.Exec(stmt, testBook1.Id)
	if err != nil {
		t.Errorf("db: delete book %d failed: %v", testBook1.Id, err)
	}

	err = ts.DeleteAuthor(testAuthor1.Id)
	is.NoErr(err)

	_, err = ts.GetAuthor(testAuthor1.Id)
	if err == nil {
		t.Errorf("expected error, author %d not deleted", testAuthor1.Id)
	}

	// check entries deleted from book_author_link
	var dest []int
	stmt = `SELECT book FROM book_author_link WHERE author=$1`
	if err := ts.db.Select(&dest, stmt, testAuthor1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from book_author_link for author %d", testAuthor1.Id)
	}
}

func TestDeleteAuthorOfExistingBook(t *testing.T) {
	err := ts.DeleteAuthor(testAuthor1.Id)
	if err == nil {
		t.Errorf("expected err: FOREIGN KEY constraint failed")
	}
}

func TestDeleteAuthorNotExists(t *testing.T) {
	err := ts.DeleteAuthor(testAuthor1.Id)
	if err == nil {
		t.Errorf("expected error: author not exists")
	}
}
