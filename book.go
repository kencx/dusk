package dusk

import (
	"dusk/validator"
	"regexp"

	"github.com/guregu/null/v5"
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

	DateStarted   null.Time `json:"-" db:"dateStarted"`
	DateCompleted null.Time `json:"-" db:"dateCompleted"`
	DateAdded     null.Time `json:"-" db:"dateAdded"`
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
		b.ISBN.ValueOrZero() != "",
		b.ISBN13.ValueOrZero() != "",
		"isbn",
		"isbn13",
		"value is missing",
	)
	if b.ISBN.ValueOrZero() != "" {
		err.Check(validator.Matches(b.ISBN.ValueOrZero(), isbnRgx), "isbn", "incorrect format")
	}
	if b.ISBN13.ValueOrZero() != "" {
		err.Check(validator.Matches(b.ISBN13.ValueOrZero(), isbnRgx), "isbn13", "incorrect format")
	}

	err.Check(b.NumOfPages >= 0, "numOfPages", "must be >= 0")
	err.Check(b.Progress >= 0, "progress", "must be >= 0")

	err.Check(b.Rating >= 0, "rating", "must be >= 0")
	err.Check(b.Rating <= 10, "rating", "must be <= 10")

	return err
}
