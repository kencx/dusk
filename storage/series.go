package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/page"
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
			return nil, fmt.Errorf("[db] failed to retrieve series %d: %w", id, err)
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
			return nil, fmt.Errorf("[db] failed to query series: %w", err)
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

func (s *Store) GetAllBooksFromSeries(id int64, f *filters.Book) (*page.Page[dusk.Book], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQueryRow
		query := buildPagedStmt(&f.Base, "book_view", `WHERE b.id IN (SELECT bookId FROM series WHERE id=$1)`)

		slog.Info("Running SQL query",
			slog.String("stmt", query),
			slog.Int64("id", id),
			slog.Int("afterId", f.AfterId),
			slog.Int("pageSize", f.Limit),
		)

		err := tx.Select(&dest, query, id, f.AfterId, f.Limit)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to retrieve books from series %d: %w", id, err)
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

func (s *Store) UpdateSeries(id int64, a *dusk.Series) (*dusk.Series, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `UPDATE series SET name=$1 WHERE id=$2`
		res, err := tx.Exec(stmt, a.Name, id)

		if err != nil {
			return nil, fmt.Errorf("[db] failed to update series %d: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to update series %d: %w", id, err)
		}
		if count == 0 {
			return nil, errors.New("[db] no series updated")
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
			return nil, fmt.Errorf("[db] failed to update series %s: %w", name, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to update series %s: %w", name, err)
		}
		if count == 0 {
			return nil, errors.New("[db] no series updated")
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
			return nil, fmt.Errorf("[db] failed to delete series %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to delete series %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("[db] series %d not removed", id)
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
		return nil, fmt.Errorf("failed to query series from book %d: %w", bookId, err)
	}
	return &series, nil
}

// Insert given series. If series already exists, return its id instead
func insertSeries(tx *sqlx.Tx, bookId int64, t string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO series (bookId, name) VALUES ($1, $2);`
	res, err := tx.Exec(stmt, bookId, t)
	if err != nil {
		return -1, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// series.name is unique
		var id int64
		stmt := `SELECT id FROM series WHERE name=$1;`
		err := tx.Get(&id, stmt, t)
		if err != nil {
			return -1, fmt.Errorf("failed to query existing series: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("failed to query existing series: %w", err)
		}
		return id, nil
	}
}

func deleteBookFromSeries(tx *sqlx.Tx, bookId, id int64) error {
	stmt := `DELETE FROM series WHERE id=$1 AND bookId=$2;`
	res, err := tx.Exec(stmt, id, bookId)
	if err != nil {
		return fmt.Errorf("failed to delete book %d from series %d: %w", bookId, id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete book %d from series %d: %w", bookId, id, err)
	}

	if count == 0 {
		return fmt.Errorf("book %d from series %d not removed", bookId, id)
	}
	return nil
}
