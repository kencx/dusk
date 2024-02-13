package storage

import (
	"database/sql"
	"dusk"
	"dusk/util"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type BookQuery struct {
	*dusk.Book
	AuthorString string          `db:"author_string"`
	TagString    util.NullString `db:"tag_string"`
}

func (s *Store) GetBook(id int64) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest BookQuery
		stmt := `SELECT b.*, GROUP_CONCAT(DISTINCT a.name) AS author_string, GROUP_CONCAT(DISTINCT t.name) AS tag_string
            FROM book b
                INNER JOIN book_author_link ba ON ba.book=b.id
                INNER JOIN author a ON ba.author=a.id
                LEFT JOIN  book_tag_link bt ON b.id=bt.book
                LEFT JOIN  tag t ON bt.tag=t.id
                LEFT JOIN  book_format_link bf ON b.id=bf.book
                LEFT JOIN  format f ON bf.format=f.id
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
		stmt := `SELECT b.*, GROUP_CONCAT(DISTINCT a.name) AS author_string, GROUP_CONCAT(DISTINCT t.name) AS tag_string
            FROM book b
                INNER JOIN book_author_link ba ON ba.book=b.id
                INNER JOIN author a ON ba.author=a.id
                LEFT JOIN  book_tag_link bt ON b.id=bt.book
                LEFT JOIN  tag t ON bt.tag=t.id
                LEFT JOIN  book_format_link bf ON b.id=bf.book
                LEFT JOIN  format f ON bf.format=f.id
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
		err = linkBookToAuthors(tx, book.ID, author_ids)
		if err != nil {
			return nil, err
		}

		if len(b.Tag) > 0 {
			tag_ids, err := insertTags(tx, b.Tag)
			if err != nil {
				return nil, err
			}
			err = linkBookToTags(tx, book.ID, tag_ids)
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

		current_authors, err := getAuthorsFromBook(tx, b.ID)
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
			if err := unlinkBookFromAuthors(tx, id, authorIDs); err != nil {
				return nil, err
			}
		}

		current_tags, err := getTagsFromBook(tx, b.ID)
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

		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}
		if err := deleteTagsWithNoBooks(tx); err != nil {
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
		// delete authors with no books
		if err := deleteAuthorsWithNoBooks(tx); err != nil {
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
		return nil, nil
	})
	return err
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b *dusk.Book) (*dusk.Book, error) {
	stmt := `INSERT INTO book (
		title,
		subtitle,
		isbn,
		isbn13,
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
		:isbn,
		:isbn13,
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
		return nil, errors.New("db: no books removed")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}

	b.ID = id
	return b, nil
}

func updateBook(tx *sqlx.Tx, id int64, b *dusk.Book) error {
	b.ID = id
	stmt := `UPDATE book
		SET
			title=:title,
			subtitle=:subtitle,
			isbn=:isbn,
			isbn13=:isbn13,
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
