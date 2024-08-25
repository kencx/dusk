package storage

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/page"
)

var (
	baseStmt = `SELECT COUNT() OVER() AS count,
			ROW_NUMBER() OVER(ORDER BY {{.SortColumn}} {{.SortDirection}}) AS rowno,
			t.*
		FROM {{.Table}} t {{.Conditional}}
	`
	pagedStmt = fmt.Sprintf(`WITH paginate AS (%s)
		SELECT * FROM paginate
		WHERE rowno > $2
		LIMIT $3;`,
		baseStmt,
	)
)

func buildBaseStmt(sortColumn, sortDirection, table, conditional string) string {
	data := map[string]interface{}{
		"SortColumn":    sortColumn,
		"SortDirection": sortDirection,
		"Table":         table,
		"Conditional":   conditional,
	}
	return Tprintf(baseStmt, data)
}

func buildPagedStmt(f *filters.Base, table, conditional string) string {
	data := map[string]interface{}{
		"SortColumn":    f.SortColumn(),
		"SortDirection": f.SortDirection(),
		"Table":         table,
		"Conditional":   conditional,
	}
	return Tprintf(pagedStmt, data)
}

func Tprintf(tmpl string, data map[string]interface{}) string {
	t := template.Must(template.New("paginationSql").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func buildSearchQuery(table string, f *filters.Search) (string, []any) {
	var (
		conditional = "WHERE $1"
		params      = []any{"1"}
	)

	if f == nil || f.Empty() {
		return buildBaseStmt("name", "ASC", table, conditional), params
	}

	// non-empty search query
	if f.Search != "" {
		conditional = fmt.Sprintf(
			`WHERE id IN (SELECT rowid FROM %[1]s_fts WHERE %[1]s_fts MATCH $1)`,
			table,
		)
		// escape search params
		params = []any{fmt.Sprintf(`"%s"`, f.Search)}
	}

	params = append(params, f.AfterId, f.Limit)
	return buildPagedStmt(&f.Base, table, conditional), params
}

func buildBookQuery(f *filters.Book) (string, []any) {
	// TODO query by:
	//   - title, subtitle
	//   - series
	var (
		conditional = "WHERE $1"
		params      = []any{"1"}
	)

	if f == nil || f.Empty() {
		return buildBaseStmt("name", "ASC", "book_view", conditional), params
	}

	// generic library search
	switch {

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
		params = []any{fmt.Sprintf(`"%s"`, f.Search.Search)}

	// ?title param
	case f.Title != "":
		conditional = `WHERE t.id IN (SELECT rowid FROM book_fts WHERE book_fts MATCH $1)`
		params = []any{fmt.Sprintf(`"%s"`, f.Title)}

	// ?author param
	case f.Author != "":
		conditional = `WHERE t.id IN (SELECT ba.book
			FROM book_author_link ba
			WHERE ba.author IN
				(SELECT rowid FROM author_fts WHERE author_fts MATCH $1))`
		params = []any{fmt.Sprintf(`"%s"`, f.Author)}

	// ?tag param
	case f.Tag != "":
		conditional = `WHERE t.id IN (SELECT bt.book
			FROM book_tag_link bt
			WHERE bt.tag IN
			(SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1))`
		params = []any{fmt.Sprintf(`"%s"`, f.Tag)}

	// no filter by default
	default:
		break
	}

	params = append(params, f.AfterId, f.Limit)
	return buildPagedStmt(&f.Base, "book_view", conditional), params
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

func newBookPage(dest []BookQueryRow, f *filters.Book) (*page.Page[dusk.Book], error) {
	// sqlx Select does not return sql.ErrNoRows
	// related issue: https://github.com/jmoiron/sqlx/issues/762#issuecomment-1062649063
	if len(dest) == 0 {
		return page.NewEmpty[dusk.Book](), nil
	}

	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, errors.New("first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(f.Limit) {
		return nil, errors.New("num of items cannot be larger than page limit")
	}

	var books []dusk.Book
	for _, row := range dest {
		row.Author = strings.Split(row.AuthorString, ",")
		row.Tag = row.TagString.Split(",")
		row.Isbn10 = row.Isbn10String.Split(",")
		row.Isbn13 = row.Isbn13String.Split(",")
		row.Formats = row.FormatString.Split(",")
		row.Series = row.SeriesString
		books = append(books, *row.Book)
	}

	result := page.New(
		int(first.Total),
		int(first.RowNo),
		int(last.RowNo),
		&f.Base,
		books,
	)
	if f.Search.Search != "" {
		result.QueryParams.Add("q", f.Search.Search)
	}
	return result, nil
}

func newAuthorPage(dest []AuthorQueryRow, f *filters.Search) (*page.Page[dusk.Author], error) {
	if len(dest) == 0 {
		return page.NewEmpty[dusk.Author](), nil
	}

	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, fmt.Errorf("first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(f.Limit) {
		return nil, fmt.Errorf("num of items cannot be larger than page limit")
	}

	var authors []dusk.Author
	for _, row := range dest {
		authors = append(authors, *row.Author)
	}

	if f == nil {
		return &page.Page[dusk.Author]{
			Info:  nil,
			Items: authors,
		}, nil
	}

	result := page.New(
		int(first.Total),
		int(first.RowNo),
		int(last.RowNo),
		&f.Base,
		authors,
	)
	if f.Search != "" {
		result.QueryParams.Add("q", f.Search)
	}
	return result, nil
}

func newTagPage(dest []TagQueryRow, f *filters.Search) (*page.Page[dusk.Tag], error) {
	if len(dest) == 0 {
		return page.NewEmpty[dusk.Tag](), nil
	}

	first := dest[0]
	last := dest[len(dest)-1]

	if first.RowNo > last.RowNo {
		return nil, fmt.Errorf("first row no cannot be larger than last row no")
	}
	if (last.RowNo - first.RowNo) > int64(f.Limit) {
		return nil, fmt.Errorf("num of items cannot be larger than page limit")
	}

	var tags []dusk.Tag
	for _, row := range dest {
		tags = append(tags, *row.Tag)
	}

	if f == nil {
		return &page.Page[dusk.Tag]{
			Info:  nil,
			Items: tags,
		}, nil
	}

	result := page.New(
		int(first.Total),
		int(first.RowNo),
		int(last.RowNo),
		&f.Base,
		tags,
	)
	if f.Search != "" {
		result.QueryParams.Add("q", f.Search)
	}
	return result, nil
}
