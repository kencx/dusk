package storage

import (
	"testing"

	"github.com/kencx/dusk"

	"github.com/matryer/is"
)

func TestGetTag(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetTag(testTag1.Id)
	is.NoErr(err)

	want := testTag1
	is.Equal(got, want)
}

func TestGetTagNotExists(t *testing.T) {
	result, err := ts.GetTag(-1)
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

// func TestGetTagWithName(t *testing.T) {
// 	got, err := ts.Tags.GetByName(testTag2.Name)
// 	checkErr(t, err)
//
// 	want := testTag2
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
// 	}
// }

func TestGetAllTags(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetAllTags()
	is.NoErr(err)

	want := allTestTags
	is.Equal(got, want)
}

func TestGetAllTagEmpty(t *testing.T) {
	defer resetDB()

	// delete all data
	if err := ts.MigrateUp(resetSchemaPath); err != nil {
		t.Errorf("failed to reset database")
	}

	got, err := ts.GetAllTags()
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

func TestGetAllBooksFromTag(t *testing.T) {
	defer resetDB()

	is := is.New(t)

	got, err := ts.GetAllBooksFromTag(testTag1.Id)
	is.NoErr(err)

	want := dusk.Books{testBook1}
	is.True(len(got) == len(want))
	for i := range want {
		if !got[i].Equal(want[i]) {
			t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
		}
	}
}

func TestCreateTag(t *testing.T) {
	is := is.New(t)
	want := &dusk.Tag{Name: "FooBar"}

	got, err := ts.CreateTag(want)
	is.NoErr(err)
	is.Equal(got.Name, want.Name)
}

func TestCreateTagDuplicates(t *testing.T) {
	is := is.New(t)
	want := testTag3

	got, err := ts.CreateTag(want)
	is.NoErr(err)
	is.Equal(got, want)

	// check for number of entries in tags
	var dest []string
	stmt := `SELECT name FROM tag WHERE name=$1`
	if err := ts.db.Select(&dest, stmt, want.Name); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(dest) != 1 {
		t.Error("more than one tag inserted")
	}
}

func TestUpdateTag(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	want := *testTag1
	want.Name = "New Tag"

	got, err := ts.UpdateTag(want.Id, &want)
	is.NoErr(err)
	is.Equal(got.Name, want.Name)
}

func TestUpdateTagExisting(t *testing.T) {
	want := *testTag1
	want.Name = testTag2.Name

	_, err := ts.UpdateTag(want.Id, &want)
	if err == nil {
		t.Errorf("expected error: unique constraint Name")
	}
}

func TestDeleteTag(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	err := ts.DeleteTag(testTag1.Id)
	is.NoErr(err)

	_, err = ts.GetTag(testTag1.Id)
	if err == nil {
		t.Errorf("expected error, tag %d not deleted", testTag1.Id)
	}

	// check entries deleted from book_tag_link
	var dest []int
	stmt := `SELECT book FROM book_tag_link WHERE tag=$1`
	if err := ts.db.Select(&dest, stmt, testTag1.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from book_tag_link for tag %d", testTag1.Id)
	}

	// check books still exist without tag
	got, err := ts.GetBook(testBook1.Id)
	is.NoErr(err)

	if len(got.Tag) != 0 {
		t.Errorf("book %d has incorrect number of tags", testBook1.Id)
	}
}

func TestDeleteTagOfBookWithRemainingTags(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	err := ts.DeleteTag(testTag3.Id)
	is.NoErr(err)

	_, err = ts.GetTag(testTag3.Id)
	if err == nil {
		t.Errorf("expected error, tag %d not deleted", testTag3.Id)
	}

	// check entries deleted from book_tag_link
	var dest []int
	stmt := `SELECT book FROM book_tag_link WHERE tag=$1`
	if err := ts.db.Select(&dest, stmt, testTag3.Id); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from book_tag_link for tag %d", testTag3.Id)
	}

	// check books still exist without tag
	got, err := ts.GetBook(testBook3.Id)
	is.NoErr(err)

	if len(got.Tag) != 1 {
		t.Errorf("book %d has incorrect number of tags", testBook3.Id)
	}
}

func TestDeleteTagNotExists(t *testing.T) {
	err := ts.DeleteTag(-1)
	if err == nil {
		t.Errorf("expected error: tag not exists")
	}
}
