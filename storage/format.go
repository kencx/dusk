package storage

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func getFormatsFromBook(tx *sqlx.Tx, bookId int64) ([]string, error) {
	var dest []struct {
		Filepath string
	}
	stmt := `SELECT i.filepath
		FROM format f
		WHERE f.bookId=$1
        ORDER BY f.id`

	if err := tx.Select(&dest, stmt, bookId); err != nil {
		return nil, err
	}

	var result []string
	for _, v := range dest {
		result = append(result, v.Filepath)
	}
	return result, nil
}

// Insert given format. If format already exists, return its id instead
func insertFormat(tx *sqlx.Tx, bookId int64, filetype, path string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO format (bookId, filetype, filepath) VALUES ($1, $2, $3);`
	res, err := tx.Exec(stmt, bookId, filetype, path)
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
func insertFormats(tx *sqlx.Tx, bookId int64, filetypes, formats []string) ([]int64, error) {
	var ids []int64

	for i, format := range formats {
		id, err := insertFormat(tx, bookId, filetypes[i], format)
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
