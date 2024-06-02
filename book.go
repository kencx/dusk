package dusk

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kennygrant/sanitize"
)

var en = language.English

type Book struct {
	Id       int64       `json:"id" db:"id"`
	Title    string      `json:"title" db:"title"`
	Subtitle null.String `json:"subtitle,omitempty" db:"subtitle"`

	// one to many
	Author []string `json:"author"`
	Tag    []string `json:"tag,omitempty"`
	Isbn10 []string `json:"isbn,omitempty"`
	Isbn13 []string `json:"isbn13,omitempty"`
	// Identifiers map[string][]string `json:"identifiers"`

	NumOfPages int `json:"num_of_pages" db:"numOfPages"`
	Rating     int `json:"rating" db:"rating"`
	Progress   int `json:"progress" db:"progress"`

	Publisher     null.String `json:"publisher" db:"publisher"`
	DatePublished null.Time   `json:"date_published" db:"datePublished"`

	Series      null.String `json:"series,omitempty" db:"series"`
	Description null.String `json:"description,omitempty" db:"description"`
	Notes       null.String `json:"notes,omitempty" db:"notes"`

	// files
	// one to many
	Formats []string    `json:"formats,omitempty"`
	Cover   null.String `json:"cover,omitempty" db:"cover"`

	DateStarted   null.Time `json:"date_started" db:"dateStarted"`
	DateCompleted null.Time `json:"date_completed" db:"dateCompleted"`
	DateAdded     null.Time `json:"date_added" db:"dateAdded"`
}

type Books []*Book

func NewBook(
	title, subtitle string,
	author, tag, formats, isbn, isbn13 []string,
	numOfPages, progress, rating int,
	publisher, series, description, notes, cover string,
	datePublished, dateStarted, dateCompleted time.Time,
) *Book {
	tcaser := cases.Title(en)
	scaser := cases.Lower(en)

	// TODO handle casing of initials
	var titleAuthor []string
	for _, a := range author {
		if a == "" {
			continue
		}
		a = strings.TrimSpace(a)
		titleAuthor = append(titleAuthor, tcaser.String(a))
	}

	var smallTag []string
	for _, a := range tag {
		if a == "" {
			continue
		}
		a = strings.TrimSpace(a)
		smallTag = append(smallTag, scaser.String(a))
	}

	b := &Book{
		Title:    tcaser.String(strings.TrimSpace(title)),
		Subtitle: null.StringFrom(tcaser.String(subtitle)),

		Author: titleAuthor,
		Tag:    smallTag,
		Isbn10: isbn,
		Isbn13: isbn13,

		NumOfPages: numOfPages,
		Progress:   progress,
		Rating:     rating,

		Publisher:     null.StringFrom(tcaser.String(publisher)),
		DatePublished: null.TimeFrom(datePublished),

		Series:      null.StringFrom(series),
		Description: null.StringFrom(description),
		Notes:       null.StringFrom(notes),

		Formats: formats,
		Cover:   null.StringFrom(cover),

		DateStarted:   null.TimeFrom(dateStarted),
		DateCompleted: null.TimeFrom(dateCompleted),
	}
	return b
}

func (b Book) SafeTitle() string {
	return sanitize.BaseName(fmt.Sprintf("%s-%d", b.Title, b.Id))
}

func (b Book) Slugify() string {
	title := strings.ReplaceAll(b.Title, ".", "")
	return sanitize.Path(fmt.Sprintf("%s-%d", title, b.Id))
}

func (b Book) Valid() validator.ErrMap {
	errMap := validator.New()

	errMap.Check(b.Title != "", "title", "value is missing")
	errMap.Check(len(b.Author) != 0, "author", "value is missing")

	for _, isbn := range b.Isbn10 {
		if isbn == "" {
			continue
		}

		ok, err := util.IsbnCheck(isbn)
		if errors.Is(err, util.ErrInvalidIsbn) {
			errMap.Add("isbn10", "invalid isbn digits")
		}
		if !ok {
			errMap.Add("isbn10", fmt.Sprintf("invalid isbn: %s", isbn))
		}
	}

	for _, isbn13 := range b.Isbn13 {
		if isbn13 == "" {
			continue
		}

		ok, err := util.IsbnCheck(isbn13)
		if errors.Is(err, util.ErrInvalidIsbn) {
			errMap.Add("isbn13", "invalid isbn digits")
		}
		if !ok {
			errMap.Add("isbn13", fmt.Sprintf("invalid isbn: %s", isbn13))
		}
	}

	errMap.Check(b.NumOfPages >= 0, "numOfPages", "must be >= 0")
	errMap.Check(b.Progress >= 0, "progress", "must be >= 0")
	errMap.Check(b.Progress >= 0, "progress", "must be <= 100")
	errMap.Check(b.Rating >= 0, "rating", "must be >= 0")
	errMap.Check(b.Rating <= 10, "rating", "must be <= 10")

	return errMap
}

func (b *Book) Equal(a *Book) bool {
	if (a == nil) && (b == nil) {
		return true
	}
	if (a != nil) && (b != nil) {
		authorEqual := reflect.DeepEqual(a.Author, b.Author)
		tagEqual := reflect.DeepEqual(a.Tag, b.Tag)
		isbn10Equal := reflect.DeepEqual(a.Isbn10, b.Isbn10)
		isbn13Equal := reflect.DeepEqual(a.Isbn13, b.Isbn13)
		formatEqual := reflect.DeepEqual(a.Formats, b.Formats)

		return (a.Title == b.Title &&
			a.Subtitle.Equal(b.Subtitle) &&
			a.NumOfPages == b.NumOfPages &&
			a.Rating == b.Rating &&
			a.Progress == b.Progress &&
			a.Publisher.Equal(b.Publisher) &&
			a.DatePublished.Equal(b.DatePublished) &&
			a.Series.Equal(b.Series) &&
			a.Description.Equal(b.Description) &&
			a.Notes.Equal(b.Notes) &&
			a.Cover.Equal(b.Cover) &&
			a.DateStarted.Equal(b.DateStarted) &&
			a.DateCompleted.Equal(b.DateCompleted) &&
			authorEqual &&
			tagEqual &&
			isbn10Equal &&
			isbn13Equal &&
			formatEqual)
	}
	return a == b
}
