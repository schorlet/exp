package test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/schorlet/exp/gtimer"
	"github.com/schorlet/exp/gtimer/storage/mem"
	"github.com/schorlet/exp/gtimer/storage/sqlite"
)

func withDB(fn func(*sqlx.DB, gtimer.TodoStore)) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer db.Close()

	var sqlStore sqlite.TodoStore
	sqlStore.MustDefine(db)

	memStore := make(mem.TodoStore)

	for _, store := range []gtimer.TodoStore{sqlStore, memStore} {
		fn(db, store)
	}
}

func TestCreateTodo(t *testing.T) {
	withDB(func(db *sqlx.DB, store gtimer.TodoStore) {
		create1, err := store.Create(db, gtimer.Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}
		if create1.Status != "active" {
			t.Fatalf("Unexpected Todo Status: %s", create1.Status)
		}

		create2, err := store.Create(db, gtimer.Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}
		if create2.ID == create1.ID {
			t.Fatalf("Unexpected Todo ID: %s", create2.ID)
		}

		todos, err := store.Read(db, gtimer.TodoFilter{})
		if err != nil {
			t.Fatalf("Unable to get Todos: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestReadTodo(t *testing.T) {
	withDB(func(db *sqlx.DB, store gtimer.TodoStore) {
		create, err := store.Create(db, gtimer.Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		todos, err := store.Read(db, gtimer.TodoFilter{ID: create.ID})
		if err != nil {
			t.Fatalf("Unable to get Todo: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todos, err = store.Read(db, gtimer.TodoFilter{ID: "0"})
		if err == nil {
			t.Fatalf("Unexpected Todo: %v", todos)
		}
		if err != gtimer.ErrNotFound {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(todos) != 0 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todos, err = store.Read(db, gtimer.TodoFilter{Status: create.Status})
		if err != nil {
			t.Fatalf("Unable to get Todo: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todos, err = store.Read(db, gtimer.TodoFilter{Status: "foo"})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(todos) != 0 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	withDB(func(db *sqlx.DB, store gtimer.TodoStore) {
		create, err := store.Create(db, gtimer.Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		create.Status = "completed"
		update, err := store.Update(db, create)
		if err != nil {
			t.Fatalf("Unable to update Todo: %v", err)
		}

		update.Status = "foo"
		update, err = store.Update(db, update)
		if err == nil {
			t.Fatal("Expected error when updating Todo")
		}

		update.ID = "0"
		create.Status = "active"
		_, err = store.Update(db, update)
		if err == nil {
			t.Fatal("Expected error when updating Todo")
		}
		if err != gtimer.ErrNotFound {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	withDB(func(db *sqlx.DB, store gtimer.TodoStore) {
		create, err := store.Create(db, gtimer.Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		err = store.Delete(db, create.ID)
		if err != nil {
			t.Fatalf("Unable to delete Todo: %v", err)
		}

		err = store.Delete(db, "0")
		if err == nil {
			t.Fatalf("Expected error when deleting Todo")
		}
		if err != gtimer.ErrNotFound {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}
