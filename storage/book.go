package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/util"

	"github.com/jmoiron/sqlx"
)

type BookRow struct {
	*dusk.Book
	AuthorString string      `db:"author_string"`
	TagString    null.String `db:"tag_string"`
	Isbn10String null.String `db:"isbn10_string"`
	Isbn13String null.String `db:"isbn13_string"`
	FormatString null.String `db:"format_string"`
	SeriesString null.String `db:"series_string"`
}

type BookQueryRow struct {
	Total int64 `db:"count"`
	RowNo int64 `db:"rowno"`
	*BookRow
}

func (s *Store) GetBook(id int64) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest BookRow
		stmt := `SELECT * FROM book_view b WHERE b.id=$1;`

		err := tx.QueryRowx(stmt, id).StructScan(&dest)
		if err == sql.ErrNoRows {
			return nil, dusk.ErrDoesNotExist
		}
		if err != nil {
			return nil, fmt.Errorf("db: retrieve book id %d failed: %w", id, err)
		}

		dest.Author = strings.Split(dest.AuthorString, ",")
		dest.Tag = dest.TagString.Split(",")
		dest.Isbn10 = dest.Isbn10String.Split(",")
		dest.Isbn13 = dest.Isbn13String.Split(",")
		dest.Formats = dest.FormatString.Split(",")
		dest.Series = null.StringFrom(dest.SeriesString.ValueOrZero())
		return dest.Book, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Book), nil
}

func (s *Store) GetAuthorsFromBook(id int64) ([]dusk.Author, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var authors []dusk.Author

		stmt := `SELECT a.*
			FROM book_author_link ba
			JOIN author a ON a.id=ba.author
			WHERE ba.book=$1
			ORDER BY a.name`

		if err := tx.Select(&authors, stmt, id); err != nil {
			return nil, err
		}
		return authors, nil
	})

	if err != nil {
		return nil, err
	}
	return i.([]dusk.Author), err
}

func (s *Store) GetAllBooks(filters *dusk.BookFilters) (*dusk.Page[dusk.Book], error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		var dest []BookQueryRow

		err := queryBooks(tx, filters, &dest)
		if err != nil {
			return nil, fmt.Errorf("db: retrieve all books with filters failed: %w", err)
		}

		// 	// sqlx Select does not return sql.ErrNoRows
		// 	// related issue: https://github.com/jmoiron/sqlx/issues/762#issuecomment-1062649063
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
				Limit:       min(int(dest[0].Total), filters.PageSize),
				TotalCount:  dest[0].Total,
				FirstRowNo:  dest[0].RowNo,
				LastRowNo:   dest[len(dest)-1].RowNo,
				QueryParams: make(url.Values),
			},
			Items: books,
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
	return i.(*dusk.Page[dusk.Book]), nil
}

func (s *Store) CreateBook(b *dusk.Book) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		book, err := insertBook(tx, b)
		if err != nil {
			return nil, err
		}

		if len(b.Author) <= 0 {
			return nil, errors.New("no authors provided")
		}

		author_ids, err := insertAuthors(tx, b.Author)
		if err != nil {
			return nil, err
		}
		err = linkBookToAuthors(tx, book.Id, author_ids)
		if err != nil {
			return nil, err
		}

		if len(b.Tag) > 0 {
			tag_ids, err := insertTags(tx, b.Tag)
			if err != nil {
				return nil, err
			}
			err = linkBookToTags(tx, book.Id, tag_ids)
			if err != nil {
				return nil, err
			}
		}

		if len(b.Isbn10) > 0 {
			_, err = insertIsbn10s(tx, book.Id, b.Isbn10)
			if err != nil {
				return nil, err
			}
		}
		if len(b.Isbn13) > 0 {
			_, err = insertIsbn13s(tx, book.Id, b.Isbn13)
			if err != nil {
				return nil, err
			}
		}

		if len(b.Formats) > 0 {
			_, err = insertFormats(tx, book.Id, b.Formats)
			if err != nil {
				return nil, err
			}
		}

		if b.Series.Valid && b.Series.ValueOrZero() != "" {
			_, err = insertSeries(tx, book.Id, b.Series.ValueOrZero())
			if err != nil {
				return nil, err
			}
		}
		return book, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Book), nil
}

func (s *Store) UpdateBook(id int64, b *dusk.Book) (*dusk.Book, error) {
	i, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		err := updateBook(tx, id, b)
		if err != nil {
			return nil, err
		}

		current_authors, err := getAuthorsFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Author)
		if !reflect.DeepEqual(current_authors, b.Author) {

			// Renaming an author should not update the same author row for other books
			// Always create a new author row, never update the original in this case
			authorIDs, err := insertAuthors(tx, b.Author)
			if err != nil {
				return nil, err
			}

			if err := linkBookToAuthors(tx, id, authorIDs); err != nil {
				return nil, err
			}

			// remove existing links for authors not in new list of ids
			if err := unlinkBookFromAuthors(tx, id, authorIDs); err != nil {
				return nil, err
			}
		}

		current_tags, err := getTagsFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Tag)
		if !reflect.DeepEqual(current_tags, b.Tag) {

			// Renaming a tag should not update the same tag row for other books
			// Always create a new tag row, never update the original in this case
			tagIDs, err := insertTags(tx, b.Tag)
			if err != nil {
				return nil, err
			}

			if err := linkBookToTags(tx, id, tagIDs); err != nil {
				return nil, err
			}
			if err := unlinkBookFromTags(tx, id, tagIDs); err != nil {
				return nil, err
			}
		}

		current_isbn10, err := getIsbn10FromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Isbn10)
		if !reflect.DeepEqual(current_isbn10, b.Isbn10) {
			if _, err = insertIsbn10s(tx, b.Id, b.Isbn10); err != nil {
				return nil, err
			}

			for _, i := range current_isbn10 {
				if err := deleteIsbn10(tx, i); err != nil {
					return nil, err
				}
			}
		}

		current_isbn13, err := getIsbn13FromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Isbn13)
		if !reflect.DeepEqual(current_isbn13, b.Isbn13) {
			if _, err = insertIsbn13s(tx, b.Id, b.Isbn13); err != nil {
				return nil, err
			}

			for _, i := range current_isbn13 {
				if err := deleteIsbn13(tx, i); err != nil {
					return nil, err
				}
			}
		}

		current_series, err := getSeriesFromBook(tx, b.Id)
		if err != nil {
			if !errors.Is(err, dusk.ErrDoesNotExist) {
				return nil, err
			}
		}

		if current_series != nil {
			if current_series.Name != b.Series.ValueOrZero() {
				if _, err = insertSeries(tx, b.Id, b.Series.ValueOrZero()); err != nil {
					return nil, err
				}

				if err := deleteBookFromSeries(tx, b.Id, current_series.Id); err != nil {
					return nil, err
				}
			}
		}

		current_formats, err := getFormatsFromBook(tx, b.Id)
		if err != nil {
			return nil, err
		}

		util.Sort(b.Formats)
		if !reflect.DeepEqual(current_formats, b.Formats) {
			if _, err = insertFormats(tx, b.Id, b.Formats); err != nil {
				return nil, err
			}
			for _, f := range current_formats {
				if err := deleteFormat(tx, f); err != nil {
					return nil, err
				}
			}
		}

		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}

		return b, nil
	})

	if err != nil {
		return nil, err
	}
	return i.(*dusk.Book), nil
}

func (s *Store) DeleteBook(id int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		if err := deleteBook(tx, id); err != nil {
			return nil, err
		}

		// delete authors with no remaining books
		// isbn10, isbn13, series and formats are
		// deleted by sqlite with CASCADE
		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (s *Store) DeleteBooks(ids []int64) error {
	_, err := Tx(s.db, func(tx *sqlx.Tx) (any, error) {
		for _, id := range ids {
			if err := deleteBook(tx, id); err != nil {
				return nil, err
			}
		}

		// delete authors with no remaining books
		// isbn10, isbn13, series and formats are
		// deleted by sqlite with CASCADE
		if err := deleteAuthorsWithNoBooks(tx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func queryBooks(tx *sqlx.Tx, filters *dusk.BookFilters, dest *[]BookQueryRow) error {
	// TODO query by:
	//   - title, subtitle
	//   - series

	var params string
	query := `WITH paginate AS (
		SELECT COUNT() OVER() AS count,
			ROW_NUMBER() OVER(ORDER BY %s %s) AS rowno,
			bv.*
		FROM book_view bv %s
	)
	SELECT * FROM paginate
	WHERE rowno > $2
	LIMIT $3;`
	query = fmt.Sprintf(query, filters.SortColumn(), filters.SortDirection(), "%s")

	switch {
	case filters == nil || filters.Empty():
		query = fmt.Sprintf(query, "WHERE $1")
		params = "1"

	// generic library search
	case filters.Search != "":
		query = fmt.Sprintf(query, `WHERE bv.id IN (
		SELECT ba.book FROM book_author_link ba
			LEFT JOIN book_tag_link bt ON ba.book=bt.book
		WHERE ba.book IN
			(SELECT rowid FROM book_fts WHERE book_fts MATCH $1)
		OR ba.author IN
			(SELECT rowid FROM author_fts WHERE author_fts MATCH $1)
		OR bt.tag IN
			(SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1))`)
		// escape params
		params = fmt.Sprintf(`"%s"`, filters.Search)

	// ?title param
	case filters.Title != "":
		query = fmt.Sprintf(query, `WHERE bv.id IN (SELECT rowid FROM book_fts WHERE book_fts MATCH $1)`)
		params = fmt.Sprintf(`"%s"`, filters.Title)

	// ?author param
	case filters.Author != "":
		query = fmt.Sprintf(query, `WHERE bv.id IN (SELECT ba.book
			FROM book_author_link ba
			WHERE ba.author IN
				(SELECT rowid FROM author_fts WHERE author_fts MATCH $1))`)
		params = fmt.Sprintf(`"%s"`, filters.Author)

	// ?tag param
	case filters.Tag != "":
		query = fmt.Sprintf(query, `WHERE bv.id IN (SELECT bt.book
			FROM book_tag_link bt
			WHERE bt.tag IN
			(SELECT rowid FROM tag_fts WHERE tag_fts MATCH $1))`)
		params = fmt.Sprintf(`"%s"`, filters.Tag)

	default:
		query = fmt.Sprintf(query, "WHERE $1")
		params = "1"
	}

	slog.Info("Running SQL query",
		slog.String("stmt", util.TrimMultiLine(query)),
		slog.Any("params", params),
		slog.Int("afterId", filters.AfterId),
		slog.Int("pageSize", filters.PageSize),
	)

	err := tx.Select(dest, query, params, filters.AfterId, filters.PageSize)
	if err != nil {
		return fmt.Errorf("db: query books failed: %w", err)
	}
	return nil
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b *dusk.Book) (*dusk.Book, error) {
	stmt := `INSERT INTO book (
		title,
		subtitle,
		numOfPages,
		progress,
		rating,
		publisher,
		datePublished,
		description,
		notes,
		cover,
		dateStarted,
		dateCompleted,
		dateAdded
	) VALUES (
		:title,
		:subtitle,
		:numOfPages,
		:progress,
		:rating,
		:publisher,
		:datePublished,
		:description,
		:notes,
		:cover,
		:dateStarted,
		:dateCompleted,
		:dateAdded);`
	res, err := tx.NamedExec(stmt, b)
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}
	if count == 0 {
		return nil, errors.New("db: no books added")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("db: insert to book table failed: %w", err)
	}

	book := *b
	book.Id = id
	return &book, nil
}

func updateBook(tx *sqlx.Tx, id int64, b *dusk.Book) error {
	b.Id = id
	stmt := `UPDATE book
		SET
			title=:title,
			subtitle=:subtitle,
			numOfPages=:numOfPages,
			progress=:progress,
			rating=:rating,
			publisher=:publisher,
			datePublished=:datePublished,
			description=:description,
			notes=:notes,
			cover=:cover,
			dateStarted=:dateStarted,
			dateCompleted=:dateCompleted,
			dateAdded=:dateAdded
		WHERE id=:id;`
	res, err := tx.NamedExec(stmt, b)

	if err != nil {
		return fmt.Errorf("db: update book %d failed: %w", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: update book %d failed: %w", id, err)
	}
	if count == 0 {
		return errors.New("db: no books updated")
	}
	return nil
}

// delete book entry from books table
func deleteBook(tx *sqlx.Tx, id int64) error {

	stmt := `DELETE from book WHERE id=$1;`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("db: delete book %d failed: %w", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete book %d failed: %w", id, err)
	}
	if count == 0 {
		return errors.New("db: no books removed")
	}
	return nil
}
