package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"

	"github.com/jmoiron/sqlx"
)

type TagQueryRow struct {
	Total int64 `db:"count"`
	RowNo int64 `db:"rowno"`
	*dusk.Tag
}

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

func (s *Store) GetAllTags(filters *dusk.SearchFilters) (*dusk.Page[dusk.Tag], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []TagQueryRow

		err := queryTags(tx, filters, &dest)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all tags with filters failed: %w", err)
		}
		if len(dest) == 0 {
			return nil, dusk.ErrNoRows
		}

		var tags []dusk.Tag
		for _, row := range dest {
			tags = append(tags, *row.Tag)
		}

		result := &dusk.Page[dusk.Tag]{
			PageInfo: &dusk.PageInfo{
				Limit:      min(int(dest[0].Total), filters.PageSize),
				TotalCount: dest[0].Total,
				FirstRowNo: dest[0].RowNo,
				LastRowNo:  dest[len(dest)-1].RowNo,
			},
			Items: tags,
		}

		// TODO
		if filters.Search != "" {
			result.QueryParams.Add("q", filters.Search)
		}
		result.QueryParams.Add("after_id", strconv.Itoa(filters.AfterId))
		result.QueryParams.Add("page_size", strconv.Itoa(filters.PageSize))
		result.QueryParams.Add("sort", filters.Sort)

		return result, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Page[dusk.Tag]), nil
}

func (s *Store) GetAllBooksFromTag(id int64, filters *dusk.BookFilters) (*dusk.Page[dusk.Book], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQueryRow

		var params string
		query := `WITH paginate AS (
			SELECT COUNT() OVER() AS count,
				ROW_NUMBER() OVER(ORDER BY %s %s) AS rowno,
				bv.*
			FROM book_view bv
			WHERE bv.id IN (SELECT book FROM book_tag_link WHERE tag=$1)
		)
		SELECT * FROM paginate
		WHERE rowno > $2
		LIMIT $3;`
		query = fmt.Sprintf(query, filters.SortColumn(), filters.SortDirection())

		slog.Info("Running SQL query",
			slog.String("stmt", query),
			slog.Int64("id", id),
			slog.Any("params", params),
			slog.Int("afterId", filters.AfterId),
			slog.Int("pageSize", filters.PageSize),
		)

		err := tx.Select(&dest, query, id, filters.AfterId, filters.PageSize)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books from tag %d failed: %w", id, err)
		}
		if len(dest) == 0 {
			return nil, dusk.ErrNoRows
		}

		var books []dusk.Book
		for _, row := range dest {
			row.Author = strings.Split(row.AuthorString, ",")
			row.Tag = row.TagString.Split(",")
			row.Isbn10 = row.Isbn10String.Split(",")
			row.Isbn13 = row.Isbn13String.Split(",")
			row.Formats = row.FormatString.Split(",")
			row.Series = null.StringFrom(row.SeriesString.ValueOrZero())
			books = append(books, *row.Book)
		}

		result := &dusk.Page[dusk.Book]{
			PageInfo: &dusk.PageInfo{
				Limit:      min(int(dest[0].Total), filters.PageSize),
				TotalCount: dest[0].Total,
				FirstRowNo: dest[0].RowNo,
				LastRowNo:  dest[len(dest)-1].RowNo,
			},
			Items: books,
		}

		// TODO
		if filters.Search != "" {
			result.QueryParams.Add("itemSearch", filters.Search)
		}
		result.QueryParams.Add("after_id", strconv.Itoa(filters.AfterId))
		result.QueryParams.Add("page_size", strconv.Itoa(filters.PageSize))
		result.QueryParams.Add("sort", filters.Sort)

		return result, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Page[dusk.Book]), nil
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

func queryTags(tx *sqlx.Tx, filters *dusk.SearchFilters, dest *[]TagQueryRow) error {
	var params string
	query := `WITH paginate AS (
		SELECT COUNT() OVER() AS count,
			ROW_NUMBER() OVER(ORDER BY %s %s) AS rowno,
			*
		FROM tag %s
	)
	SELECT * FROM paginate
	WHERE rowno > $2
	LIMIT $3;`
	query = fmt.Sprintf(query, filters.SortColumn(), filters.SortDirection(), "%s")

	switch {
	case filters == nil || filters.Empty():
		query = fmt.Sprintf(query, "WHERE $1")
		params = "1"

	case filters.Search != "":
		query = fmt.Sprintf(query, `WHERE id IN (SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1)`)
		// escape params
		params = fmt.Sprintf(`"%s"`, filters.Search)

	default:
		query = fmt.Sprintf(query, "WHERE $1")
		params = "1"
	}

	slog.Info("Running SQL query",
		slog.String("stmt", query),
		slog.Any("params", params),
		slog.Int("afterId", filters.AfterId),
		slog.Int("pageSize", filters.PageSize),
	)

	err := tx.Select(dest, query, params, filters.AfterId, filters.PageSize)
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
