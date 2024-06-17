package storage

import (
	"fmt"
	"strings"

	"github.com/kencx/dusk"
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
)

func buildPagedStmt(table string, filters *dusk.Filters, conditional string) string {
	return fmt.Sprintf(
		pageStmt,
		filters.SortColumn(),
		filters.SortDirection(),
		table,
		conditional,
	)
}

func buildPagedSearchQuery(table string, filters *dusk.SearchFilters) (string, string) {
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

	query := buildPagedStmt(table, &filters.Filters, conditional)
	return query, params
}

func buildPagedBookQuery(filters *dusk.BookFilters) (string, string) {
	// TODO query by:
	//   - title, subtitle
	//   - series
	var (
		params, conditional string
	)

	switch {
	case filters == nil || filters.Empty():
		conditional = "WHERE $1"
		params = "1"

	// generic library search
	case filters.Search != "":
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
		params = fmt.Sprintf(`"%s"`, filters.Search)

	// ?title param
	case filters.Title != "":
		conditional = `WHERE t.id IN (SELECT rowid FROM book_fts WHERE book_fts MATCH $1)`
		params = fmt.Sprintf(`"%s"`, filters.Title)

	// ?author param
	case filters.Author != "":
		conditional = `WHERE t.id IN (SELECT ba.book
			FROM book_author_link ba
			WHERE ba.author IN
				(SELECT rowid FROM author_fts WHERE author_fts MATCH $1))`
		params = fmt.Sprintf(`"%s"`, filters.Author)

	// ?tag param
	case filters.Tag != "":
		conditional = `WHERE t.id IN (SELECT bt.book
			FROM book_tag_link bt
			WHERE bt.tag IN
			(SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1))`
		params = fmt.Sprintf(`"%s"`, filters.Tag)

	default:
		conditional = "WHERE $1"
		params = "1"
	}

	query := buildPagedStmt("book_view", &filters.Filters, conditional)
	return query, params
}

type RowMetadata struct {
	Total int64 `db:"count"`
	RowNo int64 `db:"rowno"`
}

type BookQueryRow struct {
	*RowMetadata
	*BookRow
}

type AuthorQueryRow struct {
	*RowMetadata
	*dusk.Author
}

type TagQueryRow struct {
	*RowMetadata
	*dusk.Tag
}

func newBookPage(dest []BookQueryRow, filters *dusk.BookFilters) (*page.Page[dusk.Book], error) {
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
	if (last.RowNo - first.RowNo) > int64(filters.Limit) {
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
		&filters.Filters,
		books,
	)
	if filters.Search != "" {
		result.QueryParams.Add("q", filters.Search)
	}
	return result, nil
}

func newAuthorPage(dest []AuthorQueryRow, filters *dusk.SearchFilters) (*page.Page[dusk.Author], error) {
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

func newTagPage(dest []TagQueryRow, filters *dusk.SearchFilters) (*page.Page[dusk.Tag], error) {
	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, fmt.Errorf("db: first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(filters.Limit) {
		return nil, fmt.Errorf("db: num of items cannot be larger than page limit")
	}

	var tags []dusk.Tag
	for _, row := range dest {
		tags = append(tags, *row.Tag)
	}

	if len(tags) == 0 {
		return nil, dusk.ErrNoRows
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
