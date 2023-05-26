package storage

import (
	"dusk"
	"reflect"
	"testing"
)

var (
	testAuthor1 = &dusk.Author{
		ID:   1,
		Name: "John Adams",
	}
	testAuthor2 = &dusk.Author{
		ID:   2,
		Name: "Alice Brown",
	}
	testAuthor3 = &dusk.Author{
		ID:   3,
		Name: "Billy Foo",
	}
	testAuthor4 = &dusk.Author{
		ID:   4,
		Name: "Carl Baz",
	}
	testAuthor5 = &dusk.Author{
		ID:   5,
		Name: "Daniel Bar",
	}
	allTestAuthors = dusk.Authors{testAuthor1, testAuthor2, testAuthor3, testAuthor4, testAuthor5}
)

func TestGetAuthor(t *testing.T) {
	got, err := ts.GetAuthor(testAuthor1.ID)
	checkErr(t, err)

	want := testAuthor1
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestGetAuthorNotExists(t *testing.T) {
	result, err := ts.GetAuthor(-1)
	if err == nil {
		t.Errorf("expected error: ErrDoesNotExist")
	}

	if err != dusk.ErrDoesNotExist {
		t.Errorf("unexpected error: %v", err)
	}

	if result != nil {
		t.Errorf("got %v, want nil", result)
	}
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
	got, err := ts.GetAllAuthors()
	checkErr(t, err)

	want := allTestAuthors

	if len(got) != len(want) {
		t.Errorf("got %d books, want %d books", len(got), len(want))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestGetAllAuthorEmpty(t *testing.T) {
	defer resetDB()

	// delete all data
	if err := ts.MigrateUp(resetSchemaPath); err != nil {
		t.Fatalf("failed to reset database")
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

func TestCreateAuthor(t *testing.T) {
	defer resetDB()
	want := &dusk.Author{Name: "FooBar"}

	got, err := ts.CreateAuthor(want)
	checkErr(t, err)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateAuthorDuplicates(t *testing.T) {
	want := testAuthor3

	got, err := ts.CreateAuthor(want)
	checkErr(t, err)

	// check returned author is want
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

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

	want := testAuthor1
	want.Name = "Sherlock Holmes"

	got, err := ts.UpdateAuthor(want.ID, want)
	checkErr(t, err)

	if got.Name != want.Name {
		t.Errorf("got %v, want %v", got.Name, want.Name)
	}
}

func TestUpdateAuthorExisting(t *testing.T) {
	want := testAuthor1
	want.Name = testAuthor2.Name

	_, err := ts.UpdateAuthor(want.ID, want)
	if err == nil {
		t.Errorf("expected error: unique constraint Name")
	}
}

func TestDeleteAuthor(t *testing.T) {
	defer resetDB()

	// delete book first to circumvent foreign key constraint
	stmt := `DELETE from book WHERE id=$1;`
	_, err := ts.db.Exec(stmt, testBook1.ID)
	if err != nil {
		t.Errorf("db: delete book %d failed: %v", testBook1.ID, err)
	}

	err = ts.DeleteAuthor(testAuthor1.ID)
	checkErr(t, err)

	_, err = ts.GetAuthor(testAuthor1.ID)
	if err == nil {
		t.Errorf("expected error, author %d not deleted", testAuthor1.ID)
	}

	// check entries deleted from book_author_link
	var dest []int
	stmt = `SELECT book FROM book_author_link WHERE author=$1`
	if err := ts.db.Select(&dest, stmt, testAuthor1.ID); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from book_author_link for author %d", testAuthor1.ID)
	}
}

func TestDeleteAuthorOfExistingBook(t *testing.T) {
	err := ts.DeleteAuthor(testAuthor1.ID)
	if err == nil {
		t.Errorf("expected err: FOREIGN KEY constraint failed")
	}
}

func TestDeleteAuthorNotExists(t *testing.T) {
	err := ts.DeleteAuthor(testAuthor1.ID)
	if err == nil {
		t.Errorf("expected error: author not exists")
	}
}
