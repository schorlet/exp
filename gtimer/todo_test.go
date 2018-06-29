package gtimer

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func withDB(fn func(*DB, TodoStore)) {
	db := MustConnect(":memory:")
	defer db.Close()
	MustCreateSchema(db)

	fn(db, TodoSqlite{})
}

func TestCreateTodo(t *testing.T) {
	withDB(func(db *DB, store TodoStore) {
		create1, err := store.Create(db, Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		create2, err := store.Create(db, Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}
		if create2.ID == create1.ID {
			t.Fatalf("Unexpected Todo ID: %s", create2.ID)
		}

		todos, err := store.Read(db, TodoFilter{})
		if err != nil {
			t.Fatalf("Unable to get Todos: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestGetTodo(t *testing.T) {
	withDB(func(db *DB, store TodoStore) {
		create, err := store.Create(db, Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		todos, err := store.Read(db, TodoFilter{ID: create.ID})
		if err != nil {
			t.Fatalf("Unable to get Todo: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todos, err = store.Read(db, TodoFilter{ID: "0"})
		if err == nil {
			t.Fatalf("Unexpected Todo: %v", todos)
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	withDB(func(db *DB, store TodoStore) {
		create, err := store.Create(db, Todo{Title: "st101"})
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
	})
}

func TestUpdateTodo(t *testing.T) {
	withDB(func(db *DB, store TodoStore) {
		create, err := store.Create(db, Todo{Title: "st101"})
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
	})
}
