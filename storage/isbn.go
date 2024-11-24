package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/dusk"
	"github.com/mattn/go-sqlite3"
)

func getIsbn10FromBook(tx *sqlx.Tx, bookId int64) ([]string, error) {
	var dest []struct {
		Isbn string
	}
	stmt := `SELECT i.isbn
		FROM isbn10 i
		WHERE i.bookId=$1
        ORDER BY i.id`

	if err := tx.Select(&dest, stmt, bookId); err != nil {
		return nil, err
	}

	var result []string
	for _, v := range dest {
		result = append(result, v.Isbn)
	}
	return result, nil
}

// Insert given isbn10. If isbn10 already exists, return its id instead
func insertIsbn10(tx *sqlx.Tx, bookId int64, i string) (int64, error) {
	stmt := `INSERT INTO isbn10 (bookId, isbn) VALUES ($1, $2);`
	res, err := tx.Exec(stmt, bookId, i)
	if err != nil {
		if sqErr, ok := err.(sqlite3.Error); ok && sqErr.Code == sqlite3.ErrConstraint && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, dusk.ErrIsbnExists
		} else {
			return -1, err
		}
	}

	count, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if count == 0 {
		return -1, dusk.ErrNoChange
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to create isbn10: %w", err)
	}
	return id, nil
}

// Insert given slice of isbn10s and returns slice of isbn10 IDs. If isbn10 already
// exists, its ID is appended to the result
func insertIsbn10s(tx *sqlx.Tx, bookId int64, isbn10s []string) ([]int64, error) {
	var ids []int64

	for _, isbn10 := range isbn10s {
		if isbn10 != "" {
			id, err := insertIsbn10(tx, bookId, isbn10)
			if err != nil {
				return nil, fmt.Errorf("failed to insert isbn10 to book %d: %w", bookId, err)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func deleteIsbn10(tx *sqlx.Tx, isbn string) error {
	stmt := `DELETE FROM isbn10 WHERE isbn=$1;`
	res, err := tx.Exec(stmt, isbn)
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

func getIsbn13FromBook(tx *sqlx.Tx, bookId int64) ([]string, error) {
	var dest []struct {
		Isbn string
	}
	stmt := `SELECT i.isbn
		FROM isbn13 i
		WHERE i.bookId=$1
        ORDER BY i.id`

	if err := tx.Select(&dest, stmt, bookId); err != nil {
		return nil, err
	}

	var result []string
	for _, v := range dest {
		result = append(result, v.Isbn)
	}
	return result, nil
}

// Insert given isbn13. If isbn13 already exists, return its id instead
func insertIsbn13(tx *sqlx.Tx, bookId int64, i string) (int64, error) {
	stmt := `INSERT INTO isbn13 (bookId, isbn) VALUES ($1, $2);`
	res, err := tx.Exec(stmt, bookId, i)
	if err != nil {
		if sqErr, ok := err.(sqlite3.Error); ok && sqErr.Code == sqlite3.ErrConstraint && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, dusk.ErrIsbnExists
		} else {
			return -1, err
		}
	}

	count, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if count == 0 {
		return -1, dusk.ErrNoChange
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to create isbn13: %w", err)
	}
	return id, nil
}

// Insert given slice of isbn13s and returns slice of isbn13 IDs. If isbn13 already
// exists, its ID is appended to the result
func insertIsbn13s(tx *sqlx.Tx, bookId int64, isbn13s []string) ([]int64, error) {
	var ids []int64

	for _, isbn13 := range isbn13s {
		if isbn13 != "" {
			id, err := insertIsbn13(tx, bookId, isbn13)
			if err != nil {
				return nil, fmt.Errorf("failed to insert isbn13 to book %d: %w", bookId, err)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func deleteIsbn13(tx *sqlx.Tx, isbn string) error {
	stmt := `DELETE FROM isbn13 WHERE isbn=$1;`
	res, err := tx.Exec(stmt, isbn)
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
