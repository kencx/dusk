package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/util"

	"github.com/jmoiron/sqlx"
)

type BookQuery struct {
	*dusk.Book
	AuthorString string      `db:"author_string"`
	TagString    null.String `db:"tag_string"`
	Isbn10String null.String `db:"isbn10_string"`
	Isbn13String null.String `db:"isbn13_string"`
	FormatString null.String `db:"format_string"`
	SeriesString null.String `db:"series_string"`
}

func (s *Store) GetBook(id int64) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest BookQuery
		stmt := `SELECT b.*,
			GROUP_CONCAT(DISTINCT a.name) AS author_string,
			GROUP_CONCAT(DISTINCT t.name) AS tag_string,
			GROUP_CONCAT(DISTINCT it.isbn) AS isbn10_string,
			GROUP_CONCAT(DISTINCT ith.isbn) AS isbn13_string,
			GROUP_CONCAT(DISTINCT f.filepath) AS format_string,
			s.Name AS series_string
            FROM book b
                INNER JOIN book_author_link ba ON ba.book=b.id
                INNER JOIN author a ON ba.author=a.id
                LEFT JOIN  book_tag_link bt ON b.id=bt.book
                LEFT JOIN  tag t ON bt.tag=t.id
                LEFT JOIN  isbn10 it ON it.bookId=b.id
                LEFT JOIN  isbn13 ith ON ith.bookId=b.id
                LEFT JOIN  format f ON f.bookId=b.id
				LEFT JOIN  series s ON s.bookId=b.id
            WHERE b.id=$1
            GROUP BY b.id;`

		err := tx.QueryRowx(stmt, id).StructScan(&dest)
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		if err != nil {
			return nil, fmt.Errorf("db: retrieve book id %d failed: %w", id, err)
		}

		dest.Author = strings.Split(dest.AuthorString, ",")
		dest.Tag = dest.TagString.Split(",")
		dest.Isbn10 = dest.Isbn10String.Split(",")
		dest.Isbn13 = dest.Isbn13String.Split(",")
		dest.Formats = dest.FormatString.Split(",")
		dest.Series = dest.SeriesString
		return dest.Book, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Book), nil
}

func (s *Store) GetAllBooks() (dusk.Books, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQuery
		stmt := `SELECT b.*,
			GROUP_CONCAT(DISTINCT a.name) AS author_string,
			GROUP_CONCAT(DISTINCT t.name) AS tag_string,
			GROUP_CONCAT(DISTINCT it.isbn) AS isbn10_string,
			GROUP_CONCAT(DISTINCT ith.isbn) AS isbn13_string,
			GROUP_CONCAT(DISTINCT f.filepath) AS format_string,
			s.Name AS series_string
            FROM book b
                INNER JOIN book_author_link ba ON ba.book=b.id
                INNER JOIN author a ON ba.author=a.id
                LEFT JOIN  book_tag_link bt ON b.id=bt.book
                LEFT JOIN  tag t ON bt.tag=t.id
                LEFT JOIN  isbn10 it ON it.bookId=b.id
                LEFT JOIN  isbn13 ith ON ith.bookId=b.id
                LEFT JOIN  format f ON f.bookId=b.id
				LEFT JOIN  series s ON s.bookId=b.id
            GROUP BY b.id
            ORDER BY b.id;`

		err := tx.Select(&dest, stmt)
		// sqlx Select does not return sql.ErrNoRows
		// related issue: https://github.com/jmoiron/sqlx/issues/762#issuecomment-1062649063
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books failed: %w", err)
		}
		if len(dest) == 0 {
			return nil, dusk.ErrNoRows
		}

		var books dusk.Books
		for _, row := range dest {
			row.Author = strings.Split(row.AuthorString, ",")
			row.Tag = row.TagString.Split(",")
			row.Isbn10 = row.Isbn10String.Split(",")
			row.Isbn13 = row.Isbn13String.Split(",")
			row.Formats = row.FormatString.Split(",")
			row.Series = row.SeriesString
			books = append(books, row.Book)
		}
		return books, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(dusk.Books), nil
}

func (s *Store) CreateBook(b *dusk.Book) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		book, err := insertBook(tx, b)
		if err != nil {
			return nil, err
		}

		if len(b.Author) <= 0 {
			return nil, errors.New("no authors provided")
		}

		author_ids, err := insertAuthors(tx, b.Author)
		if err != nil {
			return nil, err
		}
		err = linkBookToAuthors(tx, book.Id, author_ids)
		if err != nil {
			return nil, err
		}

		if len(b.Tag) > 0 {
			tag_ids, err := insertTags(tx, b.Tag)
			if err != nil {
				return nil, err
			}
			err = linkBookToTags(tx, book.Id, tag_ids)
			if err != nil {
				return nil, err
			}
		}

		if len(b.Isbn10) > 0 {
			_, err = insertIsbn10s(tx, book.Id, b.Isbn10)
			if err != nil {
				return nil, err
			}
		}
		if len(b.Isbn13) > 0 {
			_, err = insertIsbn13s(tx, book.Id, b.Isbn13)
			if err != nil {
				return nil, err
			}
		}

		if len(b.Formats) > 0 {
			_, err = insertFormats(tx, book.Id, []string{}, b.Formats)
			if err != nil {
				return nil, err
			}
		}

		if b.Series.Valid && b.Series.ValueOrZero() != "" {
			_, err = insertSeries(tx, b.Id, b.Series.String)
			if err != nil {
				return nil, err
			}
		}
		return book, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Book), nil
}

func (s *Store) UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		err := updateBook(tx, id, b)
		if err != nil {
			return nil, err
		}

		current_authors, err := getAuthorsFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Author)
		if !reflect.DeepEqual(current_authors, b.Author) {

			// Renaming an author should not update the same author row for other books
			// Always create a new author row, never update the original in this case
			authorIDs, err := insertAuthors(tx, b.Author)
			if err != nil {
				return nil, err
			}

			if err := linkBookToAuthors(tx, id, authorIDs); err != nil {
				return nil, err
			}

			// remove existing links for authors not in new list of ids
			if err := unlinkBookFromAuthors(tx, id, authorIDs); err != nil {
				return nil, err
			}
		}

		current_tags, err := getTagsFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Tag)
		if !reflect.DeepEqual(current_tags, b.Tag) {

			// Renaming a tag should not update the same tag row for other books
			// Always create a new tag row, never update the original in this case
			tagIDs, err := insertTags(tx, b.Tag)
			if err != nil {
				return nil, err
			}

			if err := linkBookToTags(tx, id, tagIDs); err != nil {
				return nil, err
			}
			if err := unlinkBookFromTags(tx, id, tagIDs); err != nil {
				return nil, err
			}
		}

		current_isbn10, err := getIsbn10FromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Isbn10)
		if !reflect.DeepEqual(current_isbn10, b.Isbn10) {
			if _, err = insertIsbn10s(tx, b.Id, b.Isbn10); err != nil {
				return nil, err
			}
		}

		current_isbn13, err := getIsbn13FromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Isbn13)
		if !reflect.DeepEqual(current_isbn13, b.Isbn13) {
			if _, err = insertIsbn13s(tx, b.Id, b.Isbn13); err != nil {
				return nil, err
			}
		}

		current_formats, err := getFormatsFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Formats)
		if !reflect.DeepEqual(current_formats, b.Formats) {
			if _, err = insertFormats(tx, b.Id, []string{}, b.Formats); err != nil {
				return nil, err
			}
		}

		current_series, err := getSeriesFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		if current_series.Name != b.Series.ValueOrZero() {
			if _, err = insertFormats(tx, b.Id, []string{}, b.Formats); err != nil {
				return nil, err
			}
		}

		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteTagsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteIsbn10WithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteIsbn13WithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteFormatsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteSeriesWithNoBooks(tx); err != nil {
			return nil, err
		}

		return b, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Book), nil
}

func (s *Store) DeleteBook(id int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		if err := deleteBook(tx, id); err != nil {
			return nil, err
		}

		// TODO do batch delete instead
		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteIsbn10WithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteIsbn13WithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteFormatsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteSeriesWithNoBooks(tx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (s *Store) DeleteBooks(ids []int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		for _, id := range ids {
			if err := deleteBook(tx, id); err != nil {
				return nil, err
			}
		}

		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteIsbn10WithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteIsbn13WithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteFormatsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteSeriesWithNoBooks(tx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b *dusk.Book) (*dusk.Book, error) {
	stmt := `INSERT INTO book (
		title,
		subtitle,
		numOfPages,
		progress,
		rating,
		publisher,
		datePublished,
		description,
		notes,
		cover,
		dateStarted,
		dateCompleted,
		dateAdded
	) VALUES (
		:title,
		:subtitle,
		:numOfPages,
		:progress,
		:rating,
		:publisher,
		:datePublished,
		:description,
		:notes,
		:cover,
		:dateStarted,
		:dateCompleted,
		:dateAdded);`
	res, err := tx.NamedExec(stmt, b)
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}
	if count == 0 {
		return nil, errors.New("db: no books added")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}

	b.Id = id
	return b, nil
}

func updateBook(tx *sqlx.Tx, id int64, b *dusk.Book) error {
	b.Id = id
	stmt := `UPDATE book
		SET
			title=:title,
			subtitle=:subtitle,
			numOfPages=:numOfPages,
			progress=:progress,
			rating=:rating,
			publisher=:publisher,
			datePublished=:datePublished,
			description=:description,
			notes=:notes,
			cover=:cover,
			dateStarted=:dateStarted,
			dateCompleted=:dateCompleted,
			dateAdded=:dateAdded
		WHERE id=:id;`
	res, err := tx.NamedExec(stmt, b)

	if err != nil {
		return fmt.Errorf("db: update book %d failed: %w", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: update book %d failed: %w", id, err)
	}
	if count == 0 {
		return errors.New("db: no books updated")
	}
	return nil
}

// delete book entry from books table
func deleteBook(tx *sqlx.Tx, id int64) error {

	stmt := `DELETE from book WHERE id=$1;`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("db: delete book %d failed: %w", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete book %d failed: %w", id, err)
	}
	if count == 0 {
		return errors.New("db: no books removed")
	}
	return nil
}
