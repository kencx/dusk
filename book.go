package dusk

import (
	"database/sql"
	"dusk/util"
	"dusk/validator"
	"regexp"

	"github.com/kennygrant/sanitize"
)

type Book struct {
	ID          int64           `json:"id" db:"id"`
	Title       string          `json:"title" db:"title"`
	Author      []string        `json:"author"`
	ISBN        string          `json:"isbn" db:"isbn"`
	NumOfPages  int             `json:"num_of_pages" db:"numOfPages"`
	Rating      int             `json:"rating" db:"rating"`
	Description util.NullString `json:"description,omitempty" db:"description"`
	Notes       util.NullString `json:"notes,omitempty" db:"notes"`
	Tag         []string        `json:"tag,omitempty"`
	Cover       string          `json:"cover,omitempty" db:"cover"`
	Formats     []string        `json:"formats,omitempty"`

	DateCompleted sql.NullTime `json:"-" db:"dateCompleted"`
	DateAdded     sql.NullTime `json:"-" db:"dateAdded"`
	DateUpdated   sql.NullTime `json:"-" db:"dateUpdated"`
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

	err.Check(b.ISBN != "", "isbn", "value is missing")
	err.Check(validator.Matches(b.ISBN, isbnRgx), "isbn", "incorrect format")

	err.Check(b.NumOfPages >= 0, "numOfPages", "must be >= 0")

	err.Check(b.Rating >= 0, "rating", "must be >= 0")
	err.Check(b.Rating <= 10, "rating", "must be <= 10")

	return err
}
