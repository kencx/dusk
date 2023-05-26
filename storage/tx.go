package storage

import (
	"github.com/jmoiron/sqlx"
)

func Tx(db *sqlx.DB, txFunc func(*sqlx.Tx) (any, error)) (any, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	res, err := txFunc(tx)
	return res, err
}
