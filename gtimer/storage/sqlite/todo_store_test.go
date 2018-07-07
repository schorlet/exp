package sqlite

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/schorlet/exp/gtimer/storage"

	_ "github.com/mattn/go-sqlite3"
)

func todoTester(fn storage.TodoTest) func(*testing.T) {
	return func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		defer db.Close()

		var store TodoStore
		store.MustDefine(db)

		fn(t, db, store)
	}
}

func TestSqlite(t *testing.T) {
	storage.TodoTestSuite(t, todoTester)
}
