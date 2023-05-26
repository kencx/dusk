package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type model string

const (
	author = model("author")
	tag    = model("tag")
)

func linkBookToAuthors(tx *sqlx.Tx, bookID int64, ids []int64) error {
	return linkBookToModels(tx, author, bookID, ids)
}

func linkBookToTags(tx *sqlx.Tx, bookID int64, ids []int64) error {
	return linkBookToModels(tx, tag, bookID, ids)
}

func linkBookToModels(tx *sqlx.Tx, model model, bookID int64, ids []int64) error {
	type value struct {
		BookID  int64 `db:"bookID"`
		ModelID int64 `db:"modelID"`
	}

	var args = []*value{}
	for _, a := range ids {
		args = append(args, &value{
			BookID:  bookID,
			ModelID: a,
		})
	}

	stmt := fmt.Sprintf(`INSERT INTO book_%[1]v_link (book, %[1]v) VALUES (:bookID, :modelID)
	ON CONFLICT DO NOTHING;`, model)
	_, err := tx.NamedExec(stmt, args)
	if err != nil {
		return fmt.Errorf("db: link book %[2]d to %[1]v %[3]d in book_%[1]v_link failed: %[4]v", model, bookID, ids, err)
	}
	return nil
}

func unlinkBookFromAuthors(tx *sqlx.Tx, bookID int64, ids []int64) error {
	return unlinkBookFromModels(tx, author, bookID, ids)
}

func unlinkBookFromTags(tx *sqlx.Tx, bookID int64, ids []int64) error {
	return unlinkBookFromModels(tx, tag, bookID, ids)
}

func unlinkBookFromModels(tx *sqlx.Tx, model model, book_id int64, ids []int64) error {
	stmt := fmt.Sprintf(`DELETE FROM book_%[1]v_link WHERE book=? AND %[1]v NOT IN (?);`, model)
	query, args, err := sqlx.In(stmt, book_id, ids)
	if err != nil {
		return fmt.Errorf("db: unlink %[1]vs %[2]v from book %[3]v in book_%[1]vs failed: %[4]v", model, ids, book_id, err)
	}
	query = tx.Rebind(query)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("db: unlink %[1]vs %[2]v from book %[3]v in book_%[1]vs failed: %[4]v", model, ids, book_id, err)
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
