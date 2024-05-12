package goodreads

import (
	"regexp"
	"strconv"
	"strings"
	"time"

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

	rating, err := strconv.Atoi(record[7])
	if err != nil {
		return nil, err
	}

	numOfPages, err := strconv.Atoi(record[11])
	if err != nil {
		return nil, err
	}

	authors := []string{record[2]}
	if record[4] != "" {
		authors = append(authors, strings.Split(record[4], ",")...)
	}

	tags := strings.Split(record[16], ",")
	// title, series := extractSeries(title)
	// tags = append(tags, "series."+series)

	isbn10, err := util.IsbnExtract(record[5])
	if err != nil {
		return nil, err
	}
	isbn13, err := util.IsbnExtract(record[6])
	if err != nil {
		return nil, err
	}

	datePublished, err := dateparse.ParseAny(record[12])
	if err != nil {
		return nil, err
	}

	var dateRead time.Time
	if record[14] != "" {
		dateRead, err = dateparse.ParseAny(record[14])
		if err != nil {
			return nil, err
		}
	}

	var dateAdded time.Time
	if record[15] != "" {
		dateAdded, err = dateparse.ParseAny(record[15])
		if err != nil {
			return nil, err
		}
	}

	b := dusk.NewBook(
		title, subtitle,
		isbn10, isbn13,
		numOfPages, 0, rating,
		record[9], "", record[21],
		"", authors,
		tags, nil, nil, datePublished, dateAdded, dateRead,
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

	rx := regexp.MustCompile(``)
	for _, match := range rx.FindAllStringSubmatch(full, -1) {
		if len(match) > 1 {

		}
	}
	return title, series
}
