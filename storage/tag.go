package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"

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
			return nil, fmt.Errorf("db: retrieve tag %d failed: %w", id, err)
		}
		return &tag, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Tag), nil
}

func (s *Store) GetAllTags(filters *dusk.SearchFilters) (dusk.Tags, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest dusk.Tags

		if filters != nil && !filters.Empty() {
			err := queryTags(tx, *filters, &dest)
			if err != nil {
				return nil, fmt.Errorf("db: retrieve all tags with filters failed: %w", err)
			}
			if len(dest) == 0 {
				return nil, dusk.ErrNoRows
			}
		} else {
			stmt := `SELECT * FROM tag;`

			err := tx.Select(&dest, stmt)
			if err != nil {
				return nil, fmt.Errorf("db: retrieve all tags failed: %w", err)
			}
			if len(dest) == 0 {
				return nil, dusk.ErrNoRows
			}
		}
		return dest, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(dusk.Tags), nil
}

func (s *Store) GetAllBooksFromTag(id int64) (dusk.Books, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookRows
		stmt := `SELECT b.*,
			GROUP_CONCAT(DISTINCT a.name) AS author_string,
			GROUP_CONCAT(DISTINCT t.name) AS tag_string,
			GROUP_CONCAT(DISTINCT it.isbn) AS isbn10_string,
			GROUP_CONCAT(DISTINCT ith.isbn) AS isbn13_string,
			GROUP_CONCAT(DISTINCT f.filepath) AS format_string,
			s.Name AS series_string
            FROM book b
                INNER JOIN book_author_link ba ON ba.book=b.id
                INNER JOIN author a ON ba.author=a.id
                LEFT JOIN  book_tag_link bt ON b.id=bt.book
                LEFT JOIN  tag t ON bt.tag=t.id
                LEFT JOIN  isbn10 it ON it.bookId=b.id
                LEFT JOIN  isbn13 ith ON ith.bookId=b.id
                LEFT JOIN  format f ON f.bookId=b.id
				LEFT JOIN  series s ON s.bookId=b.id
            WHERE b.id IN (SELECT book FROM book_tag_link WHERE tag=$1)
			GROUP BY b.id
			ORDER BY b.id;`

		err := tx.Select(&dest, stmt, id)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from tag %d failed: %w", id, err)
		}
		if len(dest) == 0 {
			return nil, dusk.ErrNoRows
		}

		var books dusk.Books
		for _, row := range dest {
			row.Author = strings.Split(row.AuthorString, ",")
			row.Tag = row.TagString.Split(",")
			row.Isbn10 = row.Isbn10String.Split(",")
			row.Isbn13 = row.Isbn13String.Split(",")
			row.Formats = row.FormatString.Split(",")
			row.Series = null.StringFrom(row.SeriesString.ValueOrZero())
			books = append(books, row.Book)
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
		t.Id = id
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
			return nil, fmt.Errorf("db: update tag %d failed: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: update tag %d failed: %w", id, err)
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
			return nil, fmt.Errorf("db: delete tag %d failed: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("db: delete tag %d failed: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("db: tag %d not removed", id)
		}
		return nil, nil
	})
	return err
}

func queryTags(tx *sqlx.Tx, filters dusk.SearchFilters, dest *dusk.Tags) error {
	var params string
	query := ` SELECT * FROM tag
		WHERE %s
	;`

	if filters.Search != "" {
		query = fmt.Sprintf(query, `id IN (SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1)`)
		// escape params
		params = fmt.Sprintf(`"%s"`, filters.Search)
	} else {
		query = fmt.Sprintf(query, "1")
	}

	slog.Debug("Running FTS query", slog.String("stmt", query), slog.Any("params", params))

	err := tx.Select(dest, query, params)
	if err != nil {
		return fmt.Errorf("db: query tags failed: %w", err)
	}
	return nil
}

// Insert given tag. If tag already exists, return its id instead
func insertTag(tx *sqlx.Tx, t string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO tag (name) VALUES ($1);`
	res, err := tx.Exec(stmt, t)
	if err != nil {
		return -1, fmt.Errorf("db: insert to tag table failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to tag table failed: %w", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// tags.name is unique
		var id int64
		stmt := `SELECT id FROM tag WHERE name=$1;`
		err := tx.Get(&id, stmt, t)
		if err != nil {
			return -1, fmt.Errorf("db: query existing tag failed: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing tag failed: %w", err)
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
