package page

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/kencx/dusk"
	"github.com/matryer/is"
)

var (
	testPage = Page[dusk.Book]{
		Info: &Info{
			Limit:      5,
			TotalCount: 12,
			FirstRowNo: 1,
			LastRowNo:  5,
			QueryParams: url.Values{
				After: []string{"5"},
			},
		},
	}
)

func TestNew(t *testing.T) {
	is := is.New(t)
	filters := &dusk.Filters{
		AfterId: 5,
		Limit:   5,
		Sort:    "title",
	}
	items := []dusk.Book{}

	got := New(30, 1, 5, filters, items)
	want := &Page[dusk.Book]{
		Info: &Info{
			Limit:      5,
			TotalCount: 30,
			FirstRowNo: 1,
			LastRowNo:  5,
			QueryParams: url.Values{
				After: []string{"5"},
				Limit: []string{"5"},
				Sort:  []string{"title"},
			},
		},
		Items: items,
	}

	is.Equal(got, want)
}

func TestNext(t *testing.T) {
	is := is.New(t)
	got := testPage
	want := fmt.Sprintf("%s=%d", After, 5)
	is.Equal(got.Next(), want)
}

func TestNextLastPage(t *testing.T) {
	is := is.New(t)
	got := testPage
	got.FirstRowNo = 7
	got.LastRowNo = 12

	want := ""
	is.Equal(got.Next(), want)
}

func TestPrevious(t *testing.T) {
	is := is.New(t)
	got := testPage
	want := fmt.Sprintf("%s=%d", After, 1)
	is.Equal(got.Previous(), want)
}

func TestPreviousFirstPage(t *testing.T) {
	is := is.New(t)
	got := testPage
	got.FirstRowNo = 1
	got.LastRowNo = 5

	want := ""
	is.Equal(got.Previous(), want)
}

func TestNumOfPages(t *testing.T) {
	is := is.New(t)
	got := testPage
	is.Equal(got.NumOfPages(), 3)
}
