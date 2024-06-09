package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"

	"github.com/jmoiron/sqlx"
)

func (s *Store) GetAuthor(id int64) (*dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var author dusk.Author
		stmt := `SELECT * FROM author WHERE id=$1;`

		err := tx.QueryRowx(stmt, id).StructScan(&author)
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		if err != nil {
			return nil, fmt.Errorf("db: retrieve author %d failed: %w", id, err)
		}

		return &author, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Author), nil
}

func (s *Store) GetAllAuthors(filters *dusk.SearchFilters) (dusk.Authors, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest dusk.Authors

		if filters != nil && !filters.Empty() {
			err := queryAuthors(tx, *filters, &dest)
			if err != nil {
				return nil, fmt.Errorf("db: retrieve all authors with filters failed: %w", err)
			}
			if len(dest) == 0 {
				return nil, dusk.ErrNoRows
			}
		} else {
			stmt := `SELECT * FROM author;`

			err := tx.Select(&dest, stmt)
			if err != nil {
				return nil, fmt.Errorf("db: retrieve all authors failed: %w", err)
			}
			if len(dest) == 0 {
				return nil, dusk.ErrNoRows
			}
		}
		return dest, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(dusk.Authors), nil
}

func (s *Store) GetAllBooksFromAuthor(id int64) (dusk.Books, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookRows
		stmt := `SELECT b.*
			FROM book_view b
			WHERE b.id IN (SELECT book FROM book_author_link WHERE author=$1);`

		err := tx.Select(&dest, stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from author %d failed: %w", id, err)
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
			row.Series = null.StringFrom(row.SeriesString.ValueOrZero())
			books = append(books, row.Book)
		}

		return books, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(dusk.Books), nil
}

func (s *Store) CreateAuthor(a *dusk.Author) (*dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		id, err := insertAuthor(tx, a.Name)
		if err != nil {
			return nil, err
		}

		a.Id = id
		return a, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Author), nil
}

func (s *Store) UpdateAuthor(id int64, a *dusk.Author) (*dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `UPDATE author SET name=$1 WHERE id=$2`
		res, err := tx.Exec(stmt, a.Name, id)

		if err != nil {
			return nil, fmt.Errorf("db: update author %d failed: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: update author %d failed: %w", id, err)
		}
		if count == 0 {
			return nil, errors.New("db: no authors updated")
		}
		return a, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Author), nil
}

// Authors with existing books cannot be deleted. This constraint is introduced to
// prevent authors from being deleted while they are still linked to existing books.
// This relationship is only one way as books can be deleted, regardless if their
// authors still exist. It should also be noted that authors with no books will be
// deleted automatically in DeleteBook().
func (s *Store) DeleteAuthor(id int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `DELETE FROM author WHERE id=$1;`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				return nil, fmt.Errorf("db: unable to delete author %d of existing book: %w", id, err)
			}
			return nil, fmt.Errorf("db: unable to delete author %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: unable to delete author %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("db: author %d not removed", id)
		}
		return nil, nil
	})
	return err
}

func queryAuthors(tx *sqlx.Tx, filters dusk.SearchFilters, dest *dusk.Authors) error {
	var params string
	query := ` SELECT * FROM author
		WHERE %s
	;`

	if filters.Search != "" {
		query = fmt.Sprintf(query, `id IN (SELECT rowid FROM author_fts WHERE author_fts MATCH $1)`)
		// escape params
		params = fmt.Sprintf(`"%s"`, filters.Search)
	} else {
		query = fmt.Sprintf(query, "1")
	}

	slog.Debug("Running FTS query", slog.String("stmt", query), slog.Any("params", params))

	err := tx.Select(dest, query, params)
	if err != nil {
		return fmt.Errorf("db: query authors failed: %w", err)
	}
	return nil
}

// Insert given author. If author already exists, return its id instead
func insertAuthor(tx *sqlx.Tx, author string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO author (name) VALUES ($1);`
	res, err := tx.Exec(stmt, author)
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %w", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// authors.name is unique
		var id int64
		stmt := `SELECT id FROM author WHERE name=$1;`
		err := tx.Get(&id, stmt, author)
		if err != nil {
			return -1, fmt.Errorf("db: query existing author failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing author failed: %w", err)
		}
		return id, nil
	}
}

// Insert given slice of author names and returns slice of author IDs. If author already
// exists, its ID is appended to the result
func insertAuthors(tx *sqlx.Tx, authors []string) ([]int64, error) {
	var ids []int64

	for _, author := range authors {
		id, err := insertAuthor(tx, author)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// delete all authors that are not linked to any existing books
func deleteAuthorsWithNoBooks(tx *sqlx.Tx) error {
	stmt := `DELETE FROM author WHERE id NOT IN
				(SELECT author FROM book_author_link);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: unable to delete author with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete author with no books: %w", err)
	}

	if count != 0 {
		slog.Debug("deleted authors", slog.Int64("count", count))
	}
	return nil
}
