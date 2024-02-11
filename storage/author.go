package storage

import (
	"database/sql"
	"dusk"
	"errors"
	"fmt"
	"log"
	"strings"

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
			return nil, fmt.Errorf("db: retrieve author %d failed: %w", id, err)
		}

		return &author, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Author), nil
}

func (s *Store) GetAllAuthors() (dusk.Authors, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var authors dusk.Authors
		stmt := `SELECT * FROM author;`

		err := tx.Select(&authors, stmt)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all authors failed: %w", err)
		}
		if len(authors) == 0 {
			return nil, dusk.ErrNoRows
		}

		return authors, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(dusk.Authors), nil
}

func (s *Store) GetAllBooksFromAuthor(id int64) (dusk.Books, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var books dusk.Books

		stmt := `SELECT b.*
            FROM book b
                INNER JOIN book_author_link ba ON ba.book=b.id
            WHERE ba.author=$1`
		err := tx.Select(&books, stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from author %d failed: %w", id, err)
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

func (s *Store) CreateAuthor(a *dusk.Author) (*dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		id, err := insertAuthor(tx, a.Name)
		if err != nil {
			return nil, err
		}

		a.ID = id
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
			return nil, fmt.Errorf("db: update author %d failed: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: update author %d failed: %w", id, err)
		}
		if count == 0 {
			return nil, errors.New("db: no authors updated")
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
				return nil, fmt.Errorf("db: unable to delete author %d of existing book: %w", id, err)
			}
			return nil, fmt.Errorf("db: unable to delete author %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: unable to delete author %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("db: author %d not removed", id)
		}
		return nil, nil
	})
	return err
}

// Insert given author. If author already exists, return its id instead
func insertAuthor(tx *sqlx.Tx, author string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO author (name) VALUES ($1);`
	res, err := tx.Exec(stmt, author)
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %w", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// authors.name is unique
		var id int64
		stmt := `SELECT id FROM author WHERE name=$1;`
		err := tx.Get(&id, stmt, author)
		if err != nil {
			return -1, fmt.Errorf("db: query existing author failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing author failed: %w", err)
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
			return nil, err
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
		return fmt.Errorf("db: unable to delete author with no books: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete author with no books: %w", err)
	}

	if count != 0 {
		log.Printf("Deleted %d authors with no existing books", count)
		return nil
	}
	return nil
}
