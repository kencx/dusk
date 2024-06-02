package storage

import (
	"fmt"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"
)

var (
	testAuthor1    = &dusk.Author{Id: 1, Name: "Author 1"}
	testAuthor2    = &dusk.Author{Id: 2, Name: "Author 2"}
	testAuthor3    = &dusk.Author{Id: 3, Name: "Author 3"}
	testAuthor4    = &dusk.Author{Id: 4, Name: "Author 4"}
	testAuthor5    = &dusk.Author{Id: 5, Name: "Author 5"}
	allTestAuthors = dusk.Authors{testAuthor1, testAuthor2, testAuthor3, testAuthor4, testAuthor5}

	testTag1    = &dusk.Tag{Id: 1, Name: "tag 1"}
	testTag2    = &dusk.Tag{Id: 2, Name: "tag 2"}
	testTag3    = &dusk.Tag{Id: 3, Name: "tag 3"}
	allTestTags = dusk.Tags{testTag1, testTag2, testTag3}

	testIsbn101 = "0441013597"
	testIsbn102 = "0141439513"
	testIsbn131 = "9781328869333"

	testSeries1   = &dusk.Series{Id: 1, Name: "series 1"}
	allTestSeries = []*dusk.Series{testSeries1}

	testFormat1 = "format 1"

	testBook1 = &dusk.Book{
		Id:     1,
		Title:  "Book 1",
		Author: []string{testAuthor1.Name},
		Tag:    []string{testTag1.Name},
	}
	testBook2 = &dusk.Book{
		Id:      2,
		Title:   "Book 2",
		Author:  []string{testAuthor2.Name},
		Isbn10:  []string{testIsbn101},
		Series:  null.StringFrom(testSeries1.Name),
		Formats: []string{testFormat1},
	}
	testBook3 = &dusk.Book{
		Id:     3,
		Title:  "Book 3",
		Author: []string{testAuthor3.Name, testAuthor4.Name, testAuthor5.Name},
		Tag:    []string{testTag2.Name, testTag3.Name},
		Isbn10: []string{testIsbn102},
	}
	testBook4 = &dusk.Book{
		Id:     4,
		Title:  "Book 4",
		Author: []string{testAuthor5.Name},
		Isbn13: []string{testIsbn131},
	}
	allTestBooks = dusk.Books{testBook1, testBook2, testBook3, testBook4}
)

var stmts = map[string]string{
	"book":   `INSERT INTO book (title) VALUES ('%s');`,
	"author": `INSERT INTO author (name) VALUES ('%s');`,
	"tag":    `INSERT INTO tag (name) VALUES ('%s');`,

	"book_author": `INSERT INTO book_author_link (book, author) VALUES (
		(SELECT id FROM book WHERE title = '%s'),
		(SELECT id FROM author WHERE name = '%s')
	);`,
	"book_tag": `INSERT INTO book_tag_link (book, tag) VALUES (
		(SELECT id FROM book WHERE title = '%s'),
		(SELECT id FROM tag WHERE name = '%s')
	);`,

	"series": `INSERT INTO series (bookId, name) VALUES (
		(SELECT id FROM book WHERE title = '%s'), '%s'
	);`,
	"isbn10": `INSERT INTO isbn10 (bookId, isbn) VALUES (
		(SELECT id FROM book WHERE title = '%s'), '%s'
	);`,
	"isbn13": `INSERT INTO isbn13 (bookId, isbn) VALUES (
		(SELECT id FROM book WHERE title = '%s'), '%s'
	);`,
	"format": `INSERT INTO format (bookId, filepath) VALUES (
		(SELECT id FROM book WHERE title = '%s'), '%s'
	);`,
}

func runStmt(s string, v ...interface{}) error {
	stmt := fmt.Sprintf(s, v...)
	_, err := ts.db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to run stmt '%s': %w", stmt, err)
	}
	return nil
}

func seedTestBook(book *dusk.Book) error {
	if err := runStmt(stmts["book"], book.Title); err != nil {
		return err
	}

	for _, author := range book.Author {
		if err := runStmt(stmts["book_author"], book.Title, author); err != nil {
			return err
		}
	}

	for _, tag := range book.Tag {
		if err := runStmt(stmts["book_tag"], book.Title, tag); err != nil {
			return err
		}
	}

	for _, i10 := range book.Isbn10 {
		if err := runStmt(stmts["isbn10"], book.Title, i10); err != nil {
			return err
		}
	}

	for _, i13 := range book.Isbn13 {
		if err := runStmt(stmts["isbn13"], book.Title, i13); err != nil {
			return err
		}
	}

	if book.Series.Valid && book.Series.ValueOrZero() != "" {
		if err := runStmt(stmts["series"], book.Title, book.Series.ValueOrZero()); err != nil {
			return err
		}
	}

	for _, formats := range book.Formats {
		if err := runStmt(stmts["format"], book.Title, formats); err != nil {
			return err
		}
	}

	return nil
}

func seedTestAuthor(author *dusk.Author) error {
	if err := runStmt(stmts["author"], author.Name); err != nil {
		return err
	}
	return nil
}

func seedTestTag(tag *dusk.Tag) error {
	if err := runStmt(stmts["tag"], tag.Name); err != nil {
		return err
	}
	return nil
}

func seedData() error {
	for _, author := range allTestAuthors {
		if err := seedTestAuthor(author); err != nil {
			return err
		}
	}

	for _, tag := range allTestTags {
		if err := seedTestTag(tag); err != nil {
			return err
		}
	}

	for _, book := range allTestBooks {
		if err := seedTestBook(book); err != nil {
			return err
		}
	}
	return nil
}
