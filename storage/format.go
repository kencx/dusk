package storage

import (
	"database/sql"
	"fmt"
	"log"

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
		return nil, fmt.Errorf("db: query formats from book failed: %w", err)
	}

	result := dest.FormatString.Split(",")
	return result, nil
}

// Insert given format. If format already exists, return its id instead
func insertFormat(tx *sqlx.Tx, bookId int64, path string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO format (bookId, filepath) VALUES ($1, $2);`
	res, err := tx.Exec(stmt, bookId, path)
	if err != nil {
		return -1, fmt.Errorf("db: insert to format table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to format table failed: %w", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// formats.filepath is unique
		var id int64
		stmt := `SELECT id FROM format WHERE filepath=$1;`
		err := tx.Get(&id, stmt, path)
		if err != nil {
			return -1, fmt.Errorf("db: query existing format failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing format failed: %w", err)
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
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// delete all formats that are not linked to any existing books
func deleteFormatsWithNoBooks(tx *sqlx.Tx) error {
	stmt := `DELETE FROM format WHERE bookId NOT IN
				(SELECT id FROM book);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: unable to delete format with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete format with no books: %w", err)
	}

	if count != 0 {
		log.Printf("Deleted %d formats with no existing books", count)
		return nil
	}
	return nil
}
