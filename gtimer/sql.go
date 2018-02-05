package gtimer

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func withTx(db *sqlx.DB, fn func(sqlx.Ext) error) (rerr error) {
	var errs multiErr
	defer func() {
		rerr = errs.orNil()
	}()

	tx, err := db.Beginx()
	if err != nil {
		return errs.append(err)
	}

	defer func() {
		if p := recover(); p != nil {
			errs.append(fmt.Errorf("%v", p))
			if err := tx.Rollback(); err != nil {
				errs.append(err)
			}
		} else if rerr != nil {
			if err := tx.Rollback(); err != nil {
				errs.append(err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				errs.append(err)
			}
		}
	}()

	if err := fn(tx); err != nil {
		return errs.append(err)
	}
	return nil
}
