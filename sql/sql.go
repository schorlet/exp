// Package sql provides RunTx func to run a function in a transaction.
package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// MustConnect connects to the datasource or panics on error.
func MustConnect(driver, datasource string) *DB {
	db := sqlx.MustConnect(driver, datasource)
	return &DB{DB: db}
}

// Tx is a wrapper around sqlx.Tx.
type Tx struct {
	*sqlx.Tx
	commit   func() error
	rollback func() error
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	if tx.commit != nil {
		return tx.commit()
	}
	return tx.Tx.Commit()
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	if tx.rollback != nil {
		return tx.rollback()
	}
	return tx.Tx.Rollback()
}

// DB is a wrapper around sqlx.DB.
type DB struct {
	*sqlx.DB
	begin func() (*Tx, error)
}

// Beginx begins a transaction and returns an *Tx instead of an *sqlx.Tx.
func (db *DB) Beginx() (*Tx, error) {
	if db.begin != nil {
		return db.begin()
	}
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}

// RunTx runs a function in a transaction.
func RunTx(db *DB, fn func(sqlx.Ext) error) (rerr error) {
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
