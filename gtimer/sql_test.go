package gtimer

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestBegin(t *testing.T) {
	db := DB{
		begin: func() (*Tx, error) {
			return nil, fmt.Errorf("begin error")
		},
	}

	err := withTx(&db, func(sqlx.Ext) error {
		return nil
	})
	if err == nil {
		t.Fatalf("Expected error")
	}
	me, ok := err.(multiErr)
	if !ok {
		t.Fatalf("Expected multiErr")
	}
	if len(me.Errors) != 1 {
		t.Fatalf("Unexpected error count: %d", len(me.Errors))
	}
}

func newTestDB(tx *Tx) *DB {
	return &DB{
		begin: func() (*Tx, error) {
			return tx, nil
		},
	}
}

func TestCommit(t *testing.T) {
	var committed bool

	tx := Tx{
		commit: func() error {
			committed = true
			return nil
		},
	}
	db := newTestDB(&tx)

	err := withTx(db, func(sqlx.Ext) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !committed {
		t.Fatalf("Expected committed")
	}
}

func TestCommitError(t *testing.T) {
	tx := Tx{
		commit: func() error {
			return fmt.Errorf("commit error")
		},
	}
	db := newTestDB(&tx)

	err := withTx(db, func(sqlx.Ext) error {
		return nil
	})
	if err == nil {
		t.Fatalf("Expected error")
	}
	me, ok := err.(multiErr)
	if !ok {
		t.Fatalf("Expected multiErr")
	}
	if len(me.Errors) != 1 {
		t.Fatalf("len(me): %d, want: 1", len(me.Errors))
	}
}

func TestRollback(t *testing.T) {
	var rollbacked bool

	tx := Tx{
		rollback: func() error {
			rollbacked = true
			return nil
		},
	}
	db := newTestDB(&tx)

	err := withTx(db, func(sqlx.Ext) error {
		return fmt.Errorf("fn error")
	})
	if err == nil {
		t.Fatalf("Expected error")
	}
	me, ok := err.(multiErr)
	if !ok {
		t.Fatalf("Expected multiErr")
	}
	if len(me.Errors) != 1 {
		t.Fatalf("len(me): %d, want: 1", len(me.Errors))
	}
	if !rollbacked {
		t.Fatalf("Expected rollbacked")
	}
}

func TestRollbackError(t *testing.T) {
	tx := Tx{
		rollback: func() error {
			return fmt.Errorf("rollback error")
		},
	}
	db := newTestDB(&tx)

	err := withTx(db, func(sqlx.Ext) error {
		return fmt.Errorf("fn error")
	})
	if err == nil {
		t.Fatalf("Expected error")
	}
	me, ok := err.(multiErr)
	if !ok {
		t.Fatalf("Expected multiErr")
	}
	if len(me.Errors) != 2 {
		t.Fatalf("len(me): %d, want: 2", len(me.Errors))
	}
}

func TestPanic(t *testing.T) {
	var rollbacked bool

	tx := Tx{
		rollback: func() error {
			rollbacked = true
			return nil
		},
	}
	db := newTestDB(&tx)

	err := withTx(db, func(sqlx.Ext) error {
		panic("panic error")
	})
	if err == nil {
		t.Fatalf("Expected error")
	}
	me, ok := err.(multiErr)
	if !ok {
		t.Fatalf("Expected multiErr")
	}
	if len(me.Errors) != 1 {
		t.Fatalf("len(me): %d, want: 1", len(me.Errors))
	}
	if !rollbacked {
		t.Fatalf("Expected rollbacked")
	}
}

func TestPanicError(t *testing.T) {
	tx := Tx{
		rollback: func() error {
			return fmt.Errorf("rollback error")
		},
	}
	db := newTestDB(&tx)

	err := withTx(db, func(sqlx.Ext) error {
		panic("panic error")
	})
	if err == nil {
		t.Fatalf("Expected error")
	}
	me, ok := err.(multiErr)
	if !ok {
		t.Fatalf("Expected multiErr")
	}
	if len(me.Errors) != 2 {
		t.Fatalf("len(me): %d, want: 2", len(me.Errors))
	}
}
