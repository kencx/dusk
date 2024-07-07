package storage

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type model string

const (
	author = model("author")
	tag    = model("tag")
)

func linkBookToAuthors(tx *sqlx.Tx, bookId int64, ids []int64) error {
	return linkBookToModels(tx, author, bookId, ids)
}

func linkBookToTags(tx *sqlx.Tx, bookId int64, ids []int64) error {
	return linkBookToModels(tx, tag, bookId, ids)
}

func linkBookToModels(tx *sqlx.Tx, model model, bookId int64, ids []int64) error {
	if len(ids) <= 0 {
		slog.Debug(fmt.Sprintf("db: no %s were linked to book", model), slog.Int64("bookId", bookId))
		return nil
	}

	type value struct {
		BookId  int64 `db:"bookId"`
		ModelId int64 `db:"modelId"`
	}

	var args = []*value{}
	for _, a := range ids {
		args = append(args, &value{
			BookId:  bookId,
			ModelId: a,
		})
	}

	stmt := fmt.Sprintf(`INSERT INTO book_%[1]v_link (book, %[1]v) VALUES (:bookId, :modelId)
	ON CONFLICT DO NOTHING;`, model)
	_, err := tx.NamedExec(stmt, args)
	if err != nil {
		return fmt.Errorf("db: link book %[2]d to %[1]v %[3]d in book_%[1]v_link failed: %[4]v", model, bookId, ids, err)
	}
	return nil
}

func unlinkBookFromAuthors(tx *sqlx.Tx, bookId int64, ids []int64) error {
	return unlinkBookFromModels(tx, author, bookId, ids)
}

func unlinkBookFromTags(tx *sqlx.Tx, bookId int64, ids []int64) error {
	return unlinkBookFromModels(tx, tag, bookId, ids)
}

func unlinkBookFromModels(tx *sqlx.Tx, model model, bookId int64, ids []int64) error {
	var (
		stmt, query string
		args        []interface{}
		err         error
	)

	if len(ids) <= 0 {
		stmt := fmt.Sprintf(`DELETE FROM book_%s_link WHERE book=?;`, model)
		query, args, err = sqlx.In(stmt, bookId)
		if err != nil {
			return fmt.Errorf("db: unlink all %[1]vs from book %[2]v in book_%[1]v_link failed: %[3]v", model, bookId, err)
		}
	} else {
		stmt = fmt.Sprintf(`DELETE FROM book_%[1]v_link WHERE book=? AND %[1]v NOT IN (?);`, model)
		query, args, err = sqlx.In(stmt, bookId, ids)
		if err != nil {
			return fmt.Errorf("db: unlink %[1]vs %[2]v from book %[3]v in book_%[1]v_link failed: %[4]v", model, ids, bookId, err)
		}
	}

	query = tx.Rebind(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("db: unlink %[1]vs from book %[2]v in book_%[1]v_link failed: %[3]v", model, bookId, err)
	}

	if len(ids) <= 0 {
		slog.Debug(fmt.Sprintf("db: all %[1]vs were unlinked from book in book_%[1]v_link", model), slog.Int64("bookId", bookId))
	}
	return nil
}

func getAuthorsFromBook(tx *sqlx.Tx, id int64) ([]string, error) {
	return getModelsFromBook(tx, id, author)
}

func getTagsFromBook(tx *sqlx.Tx, id int64) ([]string, error) {
	return getModelsFromBook(tx, id, tag)
}

// get list of names from book id
func getModelsFromBook(tx *sqlx.Tx, id int64, model model) ([]string, error) {
	var dest []struct {
		Name string
	}
	stmt := fmt.Sprintf(`SELECT a.name
		FROM book_%[1]v_link ba
		JOIN %[1]v a ON a.id=ba.%[1]v
		WHERE ba.book=$1
        ORDER BY a.name`, model)

	if err := tx.Select(&dest, stmt, id); err != nil {
		return nil, err
	}

	var result []string
	for _, v := range dest {
		result = append(result, v.Name)
	}
	return result, nil
}
