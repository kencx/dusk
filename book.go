package dusk

import (
	"database/sql"
	"dusk/util"
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
