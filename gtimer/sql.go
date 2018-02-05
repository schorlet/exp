package gtimer

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func withTx(db *sqlx.DB, fn func(sqlx.Ext) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}
	return nil
}
