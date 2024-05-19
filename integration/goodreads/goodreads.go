package goodreads

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/util"
)

var headers = []string{
	"Book Id",
	"Title",
	"Author",
	"Author l-f",
	"Additional Authors",
	"ISBN",
	"ISBN13",
	"My Rating",
	"Average Rating",
	"Publisher",
	"Binding",
	"Number of Pages",
	"Year Published",
	"Original Publication Year",
	"Date Read",
	"Date Added",
	"Bookshelves",
	"Bookshelves with positions",
	"Exclusive Shelf",
	"My Review",
	"Spoiler",
	"Private Notes",
	"Read Count",
	"Owned Copies",
}

func RecordToBook(record []string) (*dusk.Book, error) {
	title, subtitle := extractSubtitle(record[1])

	authors := []string{record[2]}
	if record[4] != "" {
		authors = append(authors, strings.Split(record[4], ",")...)
	}

	tags := strings.Split(record[16], ",")
	// title, series := extractSeries(title)
	// tags = append(tags, "series."+series)

	isbn10, _ := util.IsbnExtract(record[5])
	isbn13, _ := util.IsbnExtract(record[6])

	rating, _ := strconv.Atoi(record[7])
	numOfPages, _ := strconv.Atoi(record[11])

	datePublished, _ := dateparse.ParseAny(record[12])
	dateRead, _ := dateparse.ParseAny(record[14])
	dateAdded, _ := dateparse.ParseAny(record[15])

	b := dusk.NewBook(
		title, subtitle,
		authors, tags, nil,
		[]string{isbn10}, []string{isbn13},
		numOfPages, 0, rating,
		record[9], "", record[21], "",
		datePublished, dateAdded, dateRead,
	)

	errMap := b.Valid()
	if len(errMap) > 0 {
		return nil, errMap
	}
	return b, nil
}

func extractSubtitle(full string) (string, string) {
	var title, subtitle string

	s := strings.Split(full, ":")
	if len(s) > 1 {
		title = s[0]
		subtitle = s[1]
	} else {
		title = full
	}
	return title, subtitle
}

// TODO extractSeries
func extractSeries(full string) (string, string) {
	var title, series string

	rx := regexp.MustCompile(`([a-zA-Z ',:]+)[(]([a-zA-Z ]+)[,][ ](#\d+)[)]$`)
	for _, match := range rx.FindAllStringSubmatch(full, -1) {
		if len(match) > 1 {

		}
	}
	return title, series
}
