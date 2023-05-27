package dusk

import (
	"database/sql"
	"dusk/util"
	"dusk/validator"
	"regexp"
)

type Book struct {
	ID         int64    `json:"id" db:"id"`
	Title      string   `json:"title" db:"title"`
	Author     []string `json:"author"`
	ISBN       string   `json:"isbn" db:"isbn"`
	NumOfPages int      `json:"num_of_pages" db:"numOfPages"`
	Rating     int      `json:"rating" db:"rating"`
	State      string   `json:"state" db:"state"`

	Description util.NullString `json:"description,omitempty"`
	Notes       util.NullString `json:"notes,omitempty"`
	Series      util.NullString `json:"series,omitempty"`
	Tag         []string        `json:"tag,omitempty"`

	DateCompleted sql.NullTime `json:"-" db:"dateCompleted"`
	DateAdded     sql.NullTime `json:"-" db:"dateAdded"`
	DateUpdated   sql.NullTime `json:"-" db:"dateUpdated"`
}

type Books []*Book

var isbnRgx = regexp.MustCompile(`[0-9]+`)

func (b *Book) Validate(v *validator.Validator) {
	v.Check(b.Title != "", "title", "value is missing")

	v.Check(len(b.Author) != 0, "author", "value is missing")

	v.Check(b.ISBN != "", "isbn", "value is missing")
	v.Check(validator.Matches(b.ISBN, isbnRgx), "isbn", "incorrect format")

	v.Check(b.NumOfPages >= 0, "numOfPages", "must be >= 0")

	v.Check(b.Rating >= 0, "rating", "must be >= 0")
	v.Check(b.Rating <= 10, "rating", "must be <= 10")
}
