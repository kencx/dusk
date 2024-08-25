package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/util"

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
			return nil, fmt.Errorf("[db] failed to retrieve author %d: %w", id, err)
		}

		return &author, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Author), nil
}

func (s *Store) GetAllAuthors(f *filters.Search) (*page.Page[dusk.Author], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []AuthorQueryRow

		err := queryAuthors(tx, f, &dest)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to query authors: %w", err)
		}

		result, err := newAuthorPage(dest, f)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*page.Page[dusk.Author]), nil
}

func (s *Store) GetAllBooksFromAuthor(id int64, f *filters.Book) (*page.Page[dusk.Book], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQueryRow
		query := buildPagedStmt(&f.Filters, "book_view", `WHERE t.id IN (SELECT book FROM book_author_link WHERE author=$1)`)

		slog.Info("Running SQL query",
			slog.String("stmt", query),
			slog.Int64("id", id),
			slog.Int("afterId", f.AfterId),
			slog.Int("pageSize", f.Limit),
		)

		err := tx.Select(&dest, query, id, f.AfterId, f.Limit)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to retrieve books from author %d: %w", id, err)
		}

		result, err := newBookPage(dest, f)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*page.Page[dusk.Book]), nil
}

func (s *Store) CreateAuthor(a *dusk.Author) (*dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		id, err := insertAuthor(tx, a.Name)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to create author: %w", err)
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
			return nil, fmt.Errorf("[db] failed to update author %d: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to update author %d: %w", id, err)
		}
		if count == 0 {
			return nil, errors.New("[db] no authors updated")
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
				return nil, fmt.Errorf("[db] failed to delete author %d of existing book: %w", id, err)
			}
			return nil, fmt.Errorf("[db] failed to delete author %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to delete author %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("[db] failed to delete author %d", id)
		}
		return nil, nil
	})
	return err
}

func queryAuthors(tx *sqlx.Tx, filters *filters.Search, dest *[]AuthorQueryRow) error {
	query, params := buildSearchQuery("author", filters)

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

// Insert given author. If author already exists, return its id instead
func insertAuthor(tx *sqlx.Tx, author string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO author (name) VALUES ($1);`
	res, err := tx.Exec(stmt, author)
	if err != nil {
		return -1, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// authors.name is unique
		var id int64
		stmt := `SELECT id FROM author WHERE name=$1;`
		err := tx.Get(&id, stmt, author)
		if err != nil {
			return -1, fmt.Errorf("failed to query existing author: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("failed to query existing author: %w", err)
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
			return nil, fmt.Errorf("failed to insert author: %w", err)
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
		return fmt.Errorf("failed to delete author with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete author with no books: %w", err)
	}

	if count != 0 {
		slog.Debug("[db] deleted authors", slog.Int64("count", count))
	}
	return nil
}
