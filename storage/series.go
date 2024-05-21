package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/dusk"
)

func (s *Store) GetSeries(id int64) (*dusk.Series, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var series dusk.Series
		stmt := `SELECT * FROM series WHERE id=$1;`

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
		stmt := `SELECT * FROM series;`

		err := tx.Select(&series, stmt)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all seriess failed: %w", err)
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
		var books dusk.Books

		stmt := `SELECT b.*
            FROM book b
                INNER JOIN series s ON s.bookId=b.id
            WHERE s.id=$1`
		err := tx.Select(&books, stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from series %d failed: %w", id, err)
		}
		if len(books) == 0 {
			return nil, dusk.ErrNoRows
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
			return nil, errors.New("db: no seriess updated")
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
			return nil, errors.New("db: no seriess updated")
		}
		return a, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Series), nil
}

// Series with existing books CAN be deleted. Their deletion will cause the series to be
// unlinked from all relevant books. This relationship goes both ways.
// A series that has no books will be deleted automatically, but a book with no series
// will not be deleted.
func (s *Store) DeleteSeries(id int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		// delete cascaded to book_series_link table
		stmt := `DELETE FROM series WHERE id=$1;`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: unable to delete series %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: unable to delete series %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("db: series %d not removed", id)
		}
		return nil, nil
	})
	return err
}

func getSeriesFromBook(tx *sqlx.Tx, bookId int64) (*dusk.Series, error) {
	var dest dusk.Series
	stmt := `SELECT s.Name
		FROM series s
		WHERE s.bookId=$1
        ORDER BY s.id`

	if err := tx.Select(&dest, stmt, bookId); err != nil {
		return nil, err
	}
	return &dest, nil
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
		// seriess.name is unique
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

// delete all series that are not linked to any existing books
func deleteSeriesWithNoBooks(tx *sqlx.Tx) error {
	stmt := `DELETE FROM series WHERE bookId NOT IN
				(SELECT id FROM book);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: unable to delete series with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete series with no books: %w", err)
	}

	if count != 0 {
		log.Printf("Deleted %d seriess with no existing books", count)
		return nil
	}
	return nil
}
