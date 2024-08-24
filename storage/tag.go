package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/page"

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
			return nil, fmt.Errorf("[db] failed to retrieve tag %d: %w", id, err)
		}
		return &tag, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Tag), nil
}

func (s *Store) GetTagsFromBook(id int64) ([]dusk.Tag, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var tags []dusk.Tag

		stmt := `SELECT a.*
			FROM book_tag_link ba
			JOIN tag a ON a.id=ba.tag
			WHERE ba.book=$1
			ORDER BY a.name`

		if err := tx.Select(&tags, stmt, id); err != nil {
			return nil, err
		}
		return tags, nil
	})

	if err != nil {
		return nil, err
	}
	return i.([]dusk.Tag), err
}

func (s *Store) GetAllTags(f *filters.Search) (*page.Page[dusk.Tag], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []TagQueryRow

		err := queryTags(tx, f, &dest)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to query tags: %w", err)
		}

		result, err := newTagPage(dest, f)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*page.Page[dusk.Tag]), nil
}

func (s *Store) GetAllBooksFromTag(id int64, f *filters.Book) (*page.Page[dusk.Book], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQueryRow
		query := buildPagedStmt(&f.Filters, "book_view", `WHERE t.id IN (SELECT book FROM book_tag_link WHERE tag=$1)`)

		slog.Info("Running SQL query",
			slog.String("stmt", query),
			slog.Int64("id", id),
			slog.Int("afterId", f.AfterId),
			slog.Int("pageSize", f.Limit),
		)

		err := tx.Select(&dest, query, id, f.AfterId, f.Limit)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to retrieve books from tag %d: %w", id, err)
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

func (s *Store) CreateTag(t *dusk.Tag) (*dusk.Tag, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		id, err := insertTag(tx, t.Name)
		if err != nil {
			return nil, fmt.Errorf("[db] failed to create tag: %w", err)
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
			return nil, fmt.Errorf("[db] failed to update tag %d: %w", id, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to update tag %d: %w", id, err)
		}
		if count == 0 {
			return nil, errors.New("[db] no tags updated")
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
			return nil, fmt.Errorf("[db] failed to delete tag %d: %w", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("[db] failed to delete tag %d: %w", id, err)
		}

		if count == 0 {
			return nil, fmt.Errorf("[db] failed to delete tag %d", id)
		}
		return nil, nil
	})
	return err
}

func queryTags(tx *sqlx.Tx, filters *filters.Search, dest *[]TagQueryRow) error {
	query, params := buildSearchQuery("tag", filters)

	slog.Info("Running SQL query",
		slog.String("stmt", query),
		slog.Any("params", params),
	)
	err := tx.Select(dest, query, params...)
	if err != nil {
		return err
	}
	return nil
}

// Insert given tag. If tag already exists, return its id instead
func insertTag(tx *sqlx.Tx, t string) (int64, error) {
	stmt := `INSERT OR IGNORE INTO tag (name) VALUES ($1);`
	res, err := tx.Exec(stmt, t)
	if err != nil {
		return -1, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// tags.name is unique
		var id int64
		stmt := `SELECT id FROM tag WHERE name=$1;`
		err := tx.Get(&id, stmt, t)
		if err != nil {
			return -1, fmt.Errorf("[db] failed to query existing tag: %w", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("[db] query existing tag failed: %w", err)
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
			return nil, fmt.Errorf("failed to insert tag: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
