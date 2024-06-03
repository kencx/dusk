package storage

import (
	"testing"

	"github.com/kencx/dusk"
	"github.com/matryer/is"
)

func TestGetSeries(t *testing.T) {
	is := is.New(t)
	got, err := ts.GetSeries(testSeries1.Id)
	is.NoErr(err)

	want := testSeries1
	is.Equal(got, want)
}

func TestGetSeriesNotExists(t *testing.T) {
	result, err := ts.GetSeries(-1)
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

func TestGetAllSeries(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	got, err := ts.GetAllSeries()
	is.NoErr(err)

	want := allTestSeries
	is.Equal(got, want)
}

func TestGetAllSeriesEmpty(t *testing.T) {
	defer resetDB()

	// delete all data
	if err := ts.MigrateUp(resetSchemaPath); err != nil {
		t.Errorf("failed to reset database")
	}

	got, err := ts.GetAllSeries()
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

func TestGetAllBooksFromSeries(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	got, err := ts.GetAllBooksFromSeries(testSeries1.Id)
	is.NoErr(err)

	want := dusk.Books{testBook2}
	is.True(len(got) == len(want))
	for i := range want {
		if !got[i].Equal(want[i]) {
			t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
		}
	}
}

func TestUpdateSeries(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	want := *testSeries1
	want.Name = "series 2"

	got, err := ts.UpdateSeries(want.Id, &want)
	is.NoErr(err)
	is.Equal(got.Name, want.Name)
}

func TestDeleteSeries(t *testing.T) {
	defer resetDB()

	is := is.New(t)
	err := ts.DeleteSeries(testSeries1.Id)
	is.NoErr(err)

	_, err = ts.GetSeries(testSeries1.Id)
	if err == nil {
		t.Errorf("expected error, series %d not deleted", testSeries1.Id)
	}

	// check books still exist without series
	got, err := ts.GetBook(testBook1.Id)
	is.NoErr(err)

	is.True(got != nil)
	is.Equal(got.Series.ValueOrZero(), "")
}
