package storage

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
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
		return -1, fmt.Errorf("db: insert %s to isbn10 table failed: %w", i, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert %s to isbn10 table failed: %w", i, err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// isbn10s.isbn is unique
		var id int64
		stmt := `SELECT id FROM isbn10 WHERE isbn=$1;`
		err := tx.Get(&id, stmt, i)
		if err != nil {
			return -1, fmt.Errorf("db: query existing isbn10 failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing isbn10 failed: %w", err)
		}
		return id, nil
	}
}

// Insert given slice of isbn10s and returns slice of isbn10 IDs. If isbn10 already
// exists, its ID is appended to the result
func insertIsbn10s(tx *sqlx.Tx, bookId int64, isbn10s []string) ([]int64, error) {
	var ids []int64

	for _, isbn10 := range isbn10s {
		if isbn10 != "" {
			id, err := insertIsbn10(tx, bookId, isbn10)
			if err != nil {
				return nil, err
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
		return fmt.Errorf("db: delete isbn10 %s failed: %w", isbn, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete isbn10 %s failed: %w", isbn, err)
	}

	if count != 0 {
		slog.Debug("deleted isbn10", slog.Int64("count", count), slog.String("isbn10", isbn))
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
		return -1, fmt.Errorf("db: insert %s to isbn13 table failed: %w", i, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert %s to isbn13 table failed: %w", i, err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// isbn13s.isbn is unique
		var id int64
		stmt := `SELECT id FROM isbn13 WHERE isbn=$1;`
		err := tx.Get(&id, stmt, i)
		if err != nil {
			return -1, fmt.Errorf("db: query existing isbn13 failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing isbn13 failed: %w", err)
		}
		return id, nil
	}
}

// Insert given slice of isbn13s and returns slice of isbn13 IDs. If isbn13 already
// exists, its ID is appended to the result
func insertIsbn13s(tx *sqlx.Tx, bookId int64, isbn13s []string) ([]int64, error) {
	var ids []int64

	for _, isbn13 := range isbn13s {
		if isbn13 != "" {
			id, err := insertIsbn13(tx, bookId, isbn13)
			if err != nil {
				return nil, err
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
		return fmt.Errorf("db: delete isbn13 %s failed: %w", isbn, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete isbn13 %s failed: %w", isbn, err)
	}

	if count != 0 {
		slog.Debug("deleted isbn13", slog.Int64("count", count), slog.String("isbn13", isbn))
	}
	return nil
}
