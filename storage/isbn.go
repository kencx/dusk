package storage

import (
	"fmt"
	"log"

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
		return -1, fmt.Errorf("db: insert to isbn10 table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to isbn10 table failed: %w", err)
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
		id, err := insertIsbn10(tx, bookId, isbn10)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// delete all isbn10s that are not linked to any existing books
func deleteIsbn10WithNoBooks(tx *sqlx.Tx) error {
	stmt := `DELETE FROM isbn10 WHERE bookId NOT IN
				(SELECT id FROM book);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: unable to delete isbn10 with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete isbn10 with no books: %w", err)
	}

	if count != 0 {
		log.Printf("Deleted %d isbn10s with no existing books", count)
		return nil
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
		return -1, fmt.Errorf("db: insert to isbn13 table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to isbn13 table failed: %w", err)
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
		id, err := insertIsbn13(tx, bookId, isbn13)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// delete all isbn13s that are not linked to any existing books
func deleteIsbn13WithNoBooks(tx *sqlx.Tx) error {
	stmt := `DELETE FROM isbn13 WHERE bookId NOT IN
				(SELECT id FROM book);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: unable to delete isbn13 with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete isbn13 with no books: %w", err)
	}

	if count != 0 {
		log.Printf("Deleted %d isbn13s with no existing books", count)
		return nil
	}
	return nil
}
