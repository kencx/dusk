package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/util"

	"github.com/jmoiron/sqlx"
)

type BookRow struct {
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
		var dest BookRow
		stmt := `SELECT * FROM book_view b WHERE b.id=$1;`

		err := tx.QueryRowx(stmt, id).StructScan(&dest)
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		if err != nil {
			return nil, fmt.Errorf("[db] failed to retrieve book id %d: %w", id, err)
		}

		// TODO test if author name has comma
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

func (s *Store) GetAuthorsFromBook(id int64) ([]dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var authors []dusk.Author
		stmt := `SELECT a.*
			FROM book_author_link ba
			JOIN author a ON a.id=ba.author
			WHERE ba.book=$1
			ORDER BY a.name`

		if err := tx.Select(&authors, stmt, id); err != nil {
			return nil, fmt.Errorf("[db] failed to retrieve authors from book id %d: %w", id, err)
		}
		return authors, nil
	})

	if err != nil {
		return nil, err
	}
	return i.([]dusk.Author), err
}

func (s *Store) GetAllBooks(f *filters.Book) (*page.Page[dusk.Book], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQueryRow

		err := queryBooks(tx, f, &dest)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to query books: %w", err)
		}

		result, err := newBookPage(dest, f)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to create new book page: %w", err)
		}
		return result, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*page.Page[dusk.Book]), nil
}

func (s *Store) CreateBook(b *dusk.Book) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		book, err := insertBook(tx, b)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to create book: %w", err)
		}

		if len(b.Author) <= 0 {
			return nil, errors.New("no authors provided")
		}

		author_ids, err := insertAuthors(tx, b.Author)
		if err != nil {
			return nil, fmt.Errorf("[db] %w", err)
		}
		err = linkBookToAuthors(tx, book.Id, author_ids)
		if err != nil {
			return nil, fmt.Errorf("[db] %w", err)
		}

		if len(b.Tag) > 0 {
			tag_ids, err := insertTags(tx, b.Tag)
			if err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
			err = linkBookToTags(tx, book.Id, tag_ids)
			if err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
		}

		if len(b.Isbn10) > 0 {
			_, err = insertIsbn10s(tx, book.Id, b.Isbn10)
			if err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
		}
		if len(b.Isbn13) > 0 {
			_, err = insertIsbn13s(tx, book.Id, b.Isbn13)
			if err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
		}

		if len(b.Formats) > 0 {
			_, err = insertFormats(tx, book.Id, b.Formats)
			if err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
		}

		if b.Series.Valid && b.Series.ValueOrZero() != "" {
			_, err = insertSeries(tx, book.Id, b.Series.ValueOrZero())
			if err != nil {
				return nil, fmt.Errorf("[db] failed to insert series for book %d: %w", b.Id, err)
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
			return nil, fmt.Errorf("[db] failed to update book %d: %w", id, err)
		}

		current_authors, err := getAuthorsFromBook(tx, b.Id)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to get author from book %d", id)
		}

		util.Sort(b.Author)
		if !reflect.DeepEqual(current_authors, b.Author) {

			// Renaming an author should not update the same author row for other books
			// Always create a new author row, never update the original in this case
			authorIDs, err := insertAuthors(tx, b.Author)
			if err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}

			if err := linkBookToAuthors(tx, id, authorIDs); err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}

			// remove existing links for authors not in new list of ids
			if err := unlinkBookFromAuthors(tx, id, authorIDs); err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
		}

		current_tags, err := getTagsFromBook(tx, b.Id)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to get tag from book %d", id)
		}

		util.Sort(b.Tag)
		if !reflect.DeepEqual(current_tags, b.Tag) {
			if len(b.Tag) <= 0 {
				if err := unlinkBookFromTags(tx, id, []int64{}); err != nil {
					return nil, fmt.Errorf("[db] %w", err)
				}
			} else {
				// Renaming a tag should not update the same tag row for other books
				// Always create a new tag row, never update the original in this case
				tagIDs, err := insertTags(tx, b.Tag)
				if err != nil {
					return nil, fmt.Errorf("[db] %w", err)
				}

				if err := linkBookToTags(tx, id, tagIDs); err != nil {
					return nil, fmt.Errorf("[db] %w", err)
				}

				if err := unlinkBookFromTags(tx, id, tagIDs); err != nil {
					return nil, fmt.Errorf("[db] %w", err)
				}
			}
		}

		current_isbn10, err := getIsbn10FromBook(tx, b.Id)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to get isbn10 from book %d", b.Id, err)
		}

		util.Sort(b.Isbn10)
		if !reflect.DeepEqual(current_isbn10, b.Isbn10) {
			if _, err = insertIsbn10s(tx, b.Id, b.Isbn10); err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}

			for _, i := range current_isbn10 {
				if err := deleteIsbn10(tx, i); err != nil {
					return nil, fmt.Errorf("[db] failed to delete isbn10: %w", err)
				}
			}
		}

		current_isbn13, err := getIsbn13FromBook(tx, b.Id)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to get isbn10 from book %d", b.Id, err)
		}

		util.Sort(b.Isbn13)
		if !reflect.DeepEqual(current_isbn13, b.Isbn13) {
			if _, err = insertIsbn13s(tx, b.Id, b.Isbn13); err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}

			for _, i := range current_isbn13 {
				if err := deleteIsbn13(tx, i); err != nil {
					return nil, fmt.Errorf("[db] failed to delete isbn13: %w", err)
				}
			}
		}

		current_series, err := getSeriesFromBook(tx, b.Id)
		if err != nil {
			if !errors.Is(err, dusk.ErrDoesNotExist) {
				return nil, fmt.Errorf("[db] failed to get series from book %d: %w", b.Id, err)
			}
		}

		if current_series != nil {
			if current_series.Name != b.Series.ValueOrZero() {
				if _, err = insertSeries(tx, b.Id, b.Series.ValueOrZero()); err != nil {
					return nil, fmt.Errorf("[db] failed to insert series for book %d: %w", b.Id, err)
				}

				if err := deleteBookFromSeries(tx, b.Id, current_series.Id); err != nil {
					return nil, fmt.Errorf("[db] failed to delete book %d from series %d: %w", b.Id, current_series.Id, err)
				}
			}
		}

		current_formats, err := getFormatsFromBook(tx, b.Id)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to retrieve formats from book %d: %w", b.Id, err)
		}

		util.Sort(b.Formats)
		if !reflect.DeepEqual(current_formats, b.Formats) {
			if _, err = insertFormats(tx, b.Id, b.Formats); err != nil {
				return nil, fmt.Errorf("[db] %w", err)
			}
			for _, f := range current_formats {
				if err := deleteFormat(tx, f); err != nil {
					return nil, fmt.Errorf("[db] failed to delete format: %w", err)
				}
			}
		}

		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, fmt.Errorf("[db] %w", err)
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
			return nil, fmt.Errorf("[db]: failed to delete book %d: %w", id, err)
		}

		// delete authors with no remaining books
		// isbn10, isbn13, series and formats are
		// deleted by sqlite with CASCADE
		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, fmt.Errorf("[db] %w", err)
		}
		return nil, nil
	})
	return err
}

func (s *Store) DeleteBooks(ids []int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		for _, id := range ids {
			if err := deleteBook(tx, id); err != nil {
				return nil, fmt.Errorf("[db]: failed to delete book %d: %w", id, err)
			}
		}

		// delete authors with no remaining books
		// isbn10, isbn13, series and formats are
		// deleted by sqlite with CASCADE
		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, fmt.Errorf("[db] %w", err)
		}
		return nil, nil
	})
	return err
}

func queryBooks(tx *sqlx.Tx, f *filters.Book, dest *[]BookQueryRow) error {
	query, params := buildBookQuery(f)

	slog.Info("Running SQL query",
		slog.String("stmt", util.TrimMultiLine(query)),
		slog.Any("params", params),
	)

	err := tx.Select(dest, query, params...)
	if err != nil {
		return err
	}
	return nil
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
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("no books added")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	book := *b
	book.Id = id
	return &book, nil
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
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no books updated")
	}
	return nil
}

// delete book entry from books table
func deleteBook(tx *sqlx.Tx, id int64) error {

	stmt := `DELETE from book WHERE id=$1;`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no books removed")
	}
	return nil
}
