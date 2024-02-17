package dusk

import (
	"errors"
	"regexp"

	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"

	"github.com/kennygrant/sanitize"
)

type Book struct {
	ID       int64       `json:"id" db:"id"`
	Title    string      `json:"title" db:"title"`
	Subtitle null.String `json:"subtitle,omitempty" db:"subtitle"`
	Author   []string    `json:"author"`

	ISBN   null.String `json:"isbn" db:"isbn"`
	ISBN13 null.String `json:"isbn13" db:"isbn13"`

	NumOfPages int `json:"num_of_pages" db:"numOfPages"`
	Progress   int `json:"progress" db:"progress"`
	Rating     int `json:"rating" db:"rating"`

	Publisher     null.String `json:"publisher" db:"publisher"`
	DatePublished null.Time   `json:"date_published" db:"datePublished"`

	Tag         []string    `json:"tag,omitempty" db:"tag"`
	Description null.String `json:"description,omitempty" db:"description"`
	Notes       null.String `json:"notes,omitempty" db:"notes"`

	// files
	Formats []string    `json:"formats,omitempty"`
	Cover   null.String `json:"cover,omitempty" db:"cover"`

	DateStarted   null.Time `json:"date_started" db:"dateStarted"`
	DateCompleted null.Time `json:"date_completed" db:"dateCompleted"`
	DateAdded     null.Time `json:"date_added" db:"dateAdded"`
}

type Books []*Book

func (b Book) SafeTitle() string {
	return sanitize.BaseName(b.Title)
}

var isbnRgx = regexp.MustCompile(`[0-9]+`)

func (b Book) Valid() validator.ErrMap {
	err := validator.New()

	err.Check(b.Title != "", "title", "value is missing")
	err.Check(len(b.Author) != 0, "author", "value is missing")

	err.EitherOr(
		b.ISBN.Valid,
		b.ISBN13.Valid,
		"isbn",
		"isbn13",
		"value is missing",
	)

	if b.ISBN.Valid {
		ok, error := util.IsbnCheck(b.ISBN.ValueOrZero())
		if errors.Is(error, util.ErrInvalidIsbn) {
			err.Add("isbn10", "invalid isbn digits")
		}
		if !ok {
			err.Add("isbn10", "invalid isbn")
		}
	}

	if b.ISBN13.Valid {
		ok, error := util.IsbnCheck(b.ISBN13.ValueOrZero())
		if errors.Is(error, util.ErrInvalidIsbn) {
			err.Add("isbn13", "invalid isbn digits")
		}
		if !ok {
			err.Add("isbn13", "invalid isbn")
		}
	}

	err.Check(b.NumOfPages >= 0, "numOfPages", "must be >= 0")
	err.Check(b.Progress >= 0, "progress", "must be >= 0")

	err.Check(b.Rating >= 0, "rating", "must be >= 0")
	err.Check(b.Rating <= 10, "rating", "must be <= 10")

	return err
}
