package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"
)

func getFormatsFromBook(tx *sqlx.Tx, bookId int64) ([]string, error) {
	var dest struct {
		FormatString null.String `db:"format_string"`
	}
	stmt := `SELECT GROUP_CONCAT(DISTINCT f.filepath) AS format_string
		FROM format f
		WHERE f.bookId=$1
        ORDER BY f.id`

	if err := tx.QueryRowx(stmt, bookId).StructScan(&dest); err != nil {
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		return nil, err
	}

	result := dest.FormatString.Split(",")
	return result, nil
}

// Insert given format. If format already exists, return its id instead
func insertFormat(tx *sqlx.Tx, bookId int64, path string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO format (bookId, filepath) VALUES ($1, $2);`
	res, err := tx.Exec(stmt, bookId, path)
	if err != nil {
		return -1, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// formats.filepath is unique
		var id int64
		stmt := `SELECT id FROM format WHERE filepath=$1;`
		err := tx.Get(&id, stmt, path)
		if err != nil {
			return -1, fmt.Errorf("[db] failed to query existing format: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("[db] failed to query existing format: %w", err)
		}
		return id, nil
	}
}

// Insert given slice of formats and returns slice of format IDs. If format already
// exists, its ID is appended to the result
func insertFormats(tx *sqlx.Tx, bookId int64, formats []string) ([]int64, error) {
	var ids []int64

	for _, format := range formats {
		id, err := insertFormat(tx, bookId, format)
		if err != nil {
			return nil, fmt.Errorf("failed to insert format to book %d: %w", bookId, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func deleteFormat(tx *sqlx.Tx, format string) error {
	stmt := `DELETE FROM format WHERE filepath=$1;`
	res, err := tx.Exec(stmt, format)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return dusk.ErrNoChange
	}
	return nil
}
