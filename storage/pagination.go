package storage

import (
	"fmt"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/page"
)

var (
	pageStmt = `WITH paginate AS (
		SELECT COUNT() OVER() AS count,
			ROW_NUMBER() OVER(ORDER BY %s %s) AS rowno,
			t.*
		FROM %s t %s
	)
	SELECT * FROM paginate
	WHERE rowno > $2
	LIMIT $3;`
	stmt = `SELECT COUNT() OVER() AS count,
			ROW_NUMBER() OVER(ORDER BY %s %s) AS rowno,
			t.*
		FROM %s t %s
	`
)

func buildPagedStmt(table string, filters *filters.Filters, conditional string) string {
	return fmt.Sprintf(
		pageStmt,
		filters.SortColumn(),
		filters.SortDirection(),
		table,
		conditional,
	)
}

func buildPagedSearchQuery(table string, filters *filters.Search) (string, string) {
	var (
		params, conditional string
	)

	switch {
	case filters == nil || filters.Empty():
		conditional = "WHERE $1"
		params = "1"
	case filters.Search != "":
		conditional = fmt.Sprintf(`WHERE id IN (SELECT rowid FROM %s_fts WHERE %s_fts MATCH $1)`, table, table)
		// escape params
		params = fmt.Sprintf(`"%s"`, filters.Search)
	default:
		conditional = "WHERE $1"
		params = "1"
	}

	var query string
	if filters == nil || filters.Empty() {
		query = fmt.Sprintf(stmt, "name", "ASC", table, conditional)
	} else {
		query = buildPagedStmt(table, &filters.Filters, conditional)
	}

	return query, params
}

func buildPagedBookQuery(f *filters.Book) (string, string) {
	// TODO query by:
	//   - title, subtitle
	//   - series
	var (
		params, conditional string
	)

	switch {
	case f == nil || f.Empty():
		conditional = "WHERE $1"
		params = "1"

	// generic library search
	case f.Search.Search != "":
		conditional = `WHERE t.id IN (
		SELECT ba.book FROM book_author_link ba
			LEFT JOIN book_tag_link bt ON ba.book=bt.book
		WHERE ba.book IN
			(SELECT rowid FROM book_fts WHERE book_fts MATCH $1)
		OR ba.author IN
			(SELECT rowid FROM author_fts WHERE author_fts MATCH $1)
		OR bt.tag IN
			(SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1))`
		// escape params
		params = fmt.Sprintf(`"%s"`, f.Search)

	// ?title param
	case f.Title != "":
		conditional = `WHERE t.id IN (SELECT rowid FROM book_fts WHERE book_fts MATCH $1)`
		params = fmt.Sprintf(`"%s"`, f.Title)

	// ?author param
	case f.Author != "":
		conditional = `WHERE t.id IN (SELECT ba.book
			FROM book_author_link ba
			WHERE ba.author IN
				(SELECT rowid FROM author_fts WHERE author_fts MATCH $1))`
		params = fmt.Sprintf(`"%s"`, f.Author)

	// ?tag param
	case f.Tag != "":
		conditional = `WHERE t.id IN (SELECT bt.book
			FROM book_tag_link bt
			WHERE bt.tag IN
			(SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1))`
		params = fmt.Sprintf(`"%s"`, f.Tag)

	default:
		conditional = "WHERE $1"
		params = "1"
	}

	query := buildPagedStmt("book_view", &f.Filters, conditional)
	return query, params
}

type RowMetadata struct {
	Total int64 `db:"count"`
	RowNo int64 `db:"rowno"`
}

type BookQueryRow struct {
	*RowMetadata
	*bookRow
}

type AuthorQueryRow struct {
	*RowMetadata
	*dusk.Author
}

type TagQueryRow struct {
	*RowMetadata
	*dusk.Tag
}

func newBookPage(dest []BookQueryRow, f *filters.Book) (*page.Page[dusk.Book], error) {
	// sqlx Select does not return sql.ErrNoRows
	// related issue: https://github.com/jmoiron/sqlx/issues/762#issuecomment-1062649063
	if len(dest) == 0 {
		return nil, dusk.ErrNoRows
	}

	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, fmt.Errorf("db: first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(f.Limit) {
		return nil, fmt.Errorf("db: num of items cannot be larger than page limit")
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

	result := page.New(
		int(first.Total),
		int(first.RowNo),
		int(last.RowNo),
		&f.Filters,
		books,
	)
	if f.Search.Search != "" {
		result.QueryParams.Add("q", f.Search.Search)
	}
	return result, nil
}

func newAuthorPage(dest []AuthorQueryRow, filters *filters.Search) (*page.Page[dusk.Author], error) {
	if len(dest) == 0 {
		return nil, dusk.ErrNoRows
	}

	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, fmt.Errorf("db: first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(filters.Limit) {
		return nil, fmt.Errorf("db: num of items cannot be larger than page limit")
	}

	var authors []dusk.Author
	for _, row := range dest {
		authors = append(authors, *row.Author)
	}

	result := page.New(
		int(first.Total),
		int(first.RowNo),
		int(last.RowNo),
		&filters.Filters,
		authors,
	)
	if filters.Search != "" {
		result.QueryParams.Add("q", filters.Search)
	}
	return result, nil
}

func newTagPage(dest []TagQueryRow, filters *filters.Search) (*page.Page[dusk.Tag], error) {
	var tags []dusk.Tag
	for _, row := range dest {
		tags = append(tags, *row.Tag)
	}

	if len(tags) == 0 {
		return nil, dusk.ErrNoRows
	}

	if filters == nil {
		return &page.Page[dusk.Tag]{
			Info:  nil,
			Items: tags,
		}, nil
	}

	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, fmt.Errorf("db: first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(filters.Limit) {
		return nil, fmt.Errorf("db: num of items cannot be larger than page limit")
	}

	result := page.New(
		int(first.Total),
		int(first.RowNo),
		int(last.RowNo),
		&filters.Filters,
		tags,
	)
	if filters.Search != "" {
		result.QueryParams.Add("q", filters.Search)
	}
	return result, nil
}
