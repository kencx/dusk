package storage

import (
	"database/sql"
	"dusk"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func (s *Store) GetTag(id int64) (*dusk.Tag, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var tag dusk.Tag
		stmt := `SELECT * FROM tag WHERE id=$1;`

		err := tx.QueryRowx(stmt, id).StructScan(&tag)
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		if err != nil {
			return nil, fmt.Errorf("db: retrieve tag %d failed: %v", id, err)
		}
		return &tag, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Tag), nil
}

func (s *Store) GetAllTags() (dusk.Tags, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var tags dusk.Tags
		stmt := `SELECT * FROM tag;`

		err := tx.Select(&tags, stmt)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all tags failed: %v", err)
		}
		if len(tags) == 0 {
			return nil, dusk.ErrNoRows
		}
		return tags, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(dusk.Tags), nil
}

func (s *Store) GetAllBooksFromTag(id int64) (dusk.Books, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var books dusk.Books

		stmt := `SELECT b.*
            FROM book b
                INNER JOIN book_tag_link ba ON ba.book=b.id
            WHERE ba.tag=$1`
		err := tx.Select(&books, stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from tag %d failed: %v", id, err)
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

func (s *Store) CreateTag(t *dusk.Tag) (*dusk.Tag, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		id, err := insertTag(tx, t.Name)
		if err != nil {
			return nil, err
		}
		t.ID = id
		return t, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Tag), nil
}

func (s *Store) UpdateTag(id int64, a *dusk.Tag) (*dusk.Tag, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		stmt := `UPDATE tag SET name=$1 WHERE id=$2`
		res, err := tx.Exec(stmt, a.Name, id)

		if err != nil {
			return nil, fmt.Errorf("db: update tag %d failed: %v", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: update tag %d failed: %v", id, err)
		}
		if count == 0 {
			return nil, errors.New("db: no tags updated")
		}
		return a, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Tag), nil
}

// Tags with existing books CAN be deleted. Their deletion will cause the tag to be
// unlinked from all relevant books. This relationship goes both ways. A tag that has no
// books will be deleted automatically.
func (s *Store) DeleteTag(id int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		// delete cascaded to book_tag_link table
		stmt := `DELETE FROM tag WHERE id=$1;`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: unable to delete tag %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: unable to delete tag %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("db: tag %d not removed", id)
		}
		return nil, nil
	})
	return err
}

// Insert given tag. If tag already exists, return its id instead
func insertTag(tx *sqlx.Tx, t string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO tag (name) VALUES ($1);`
	res, err := tx.Exec(stmt, t)
	if err != nil {
		return -1, fmt.Errorf("db: insert to tag table failed: %v", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to tag table failed: %v", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// tags.name is unique
		var id int64
		stmt := `SELECT id FROM tag WHERE name=$1;`
		err := tx.Get(&id, stmt, t)
		if err != nil {
			return -1, fmt.Errorf("db: query existing tag failed: %v", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing tag failed: %v", err)
		}
		return id, nil
	}
}

// Insert given slice of tags and returns slice of tag IDs. If tag already
// exists, its ID is appended to the result
func insertTags(tx *sqlx.Tx, tags []string) ([]int64, error) {
	var ids []int64

	for _, tag := range tags {
		id, err := insertTag(tx, tag)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// delete all tags that are not linked to any existing books
func deleteTagsWithNoBooks(tx *sqlx.Tx) error {
	stmt := `DELETE FROM tag WHERE id NOT IN
				(SELECT tag FROM book_tag_link);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: unable to delete tag with no books: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete tag with no books: %v", err)
	}

	if count != 0 {
		log.Printf("Deleted %d tags with no existing books", count)
		return nil
	}
	return nil
}
