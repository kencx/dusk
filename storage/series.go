package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"
)

func (s *Store) GetSeries(id int64) (*dusk.Series, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var series dusk.Series
		stmt := `SELECT id, name FROM series WHERE id=$1;`

		err := tx.QueryRowx(stmt, id).StructScan(&series)
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		if err != nil {
			return nil, fmt.Errorf("db: retrieve series %d failed: %w", id, err)
		}
		return &series, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Series), nil
}

func (s *Store) GetAllSeries() ([]*dusk.Series, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var series []*dusk.Series
		stmt := `SELECT id, name FROM series;`

		err := tx.Select(&series, stmt)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all series failed: %w", err)
		}
		if len(series) == 0 {
			return nil, dusk.ErrNoRows
		}
		return series, nil
	})

	if err != nil {
		return nil, err
	}
	return i.([]*dusk.Series), nil
}

func (s *Store) GetAllBooksFromSeries(id int64) (dusk.Books, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookRows
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
            WHERE b.id IN (SELECT bookId FROM series WHERE id=$1)
			GROUP BY b.id
			ORDER BY b.id;`

		err := tx.Select(&dest, stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from series %d failed: %w", id, err)
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

func (s *Store) UpdateSeries(id int64, a *dusk.Series) (*dusk.Series, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `UPDATE series SET name=$1 WHERE id=$2`
		res, err := tx.Exec(stmt, a.Name, id)

		if err != nil {
			return nil, fmt.Errorf("db: update series %d failed: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: update series %d failed: %w", id, err)
		}
		if count == 0 {
			return nil, errors.New("db: no series updated")
		}
		return a, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Series), nil
}

func (s *Store) UpdateSeriesByName(name string, a *dusk.Series) (*dusk.Series, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `UPDATE series SET name=$1 WHERE name=$2`
		res, err := tx.Exec(stmt, a.Name, name)

		if err != nil {
			return nil, fmt.Errorf("db: update series %s failed: %w", name, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: update series %s failed: %w", name, err)
		}
		if count == 0 {
			return nil, errors.New("db: no series updated")
		}
		return a, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Series), nil
}

// Series with existing books CAN be deleted. Their deletion will cause the series to be
// unlinked from all relevant books.
func (s *Store) DeleteSeries(id int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `DELETE FROM series WHERE id=$1;`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: delete series %d failed: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: delete series %d failed: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("db: series %d not removed", id)
		}
		return nil, nil
	})
	return err
}

func getSeriesFromBook(tx *sqlx.Tx, bookId int64) (*dusk.Series, error) {
	var series dusk.Series
	stmt := `SELECT id, name
		FROM series
		WHERE bookId=$1
        ORDER BY id`

	if err := tx.QueryRowx(stmt, bookId).StructScan(&series); err != nil {
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		return nil, fmt.Errorf("db: query series from book failed: %w", err)
	}
	return &series, nil
}

// Insert given series. If series already exists, return its id instead
func insertSeries(tx *sqlx.Tx, bookId int64, t string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO series (bookId, name) VALUES ($1, $2);`
	res, err := tx.Exec(stmt, bookId, t)
	if err != nil {
		return -1, fmt.Errorf("db: insert to series table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to series table failed: %w", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// series.name is unique
		var id int64
		stmt := `SELECT id FROM series WHERE name=$1;`
		err := tx.Get(&id, stmt, t)
		if err != nil {
			return -1, fmt.Errorf("db: query existing series failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing series failed: %w", err)
		}
		return id, nil
	}
}

func deleteBookFromSeries(tx *sqlx.Tx, bookId, id int64) error {
	stmt := `DELETE FROM series WHERE id=$1 AND bookId=$2;`
	res, err := tx.Exec(stmt, id, bookId)
	if err != nil {
		return fmt.Errorf("db: delete book %d from series %d failed: %w", bookId, id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete book %d from series %d failed: %w", bookId, id, err)
	}

	if count == 0 {
		return fmt.Errorf("db: book %d from series %d not removed", bookId, id)
	}
	return nil
}
