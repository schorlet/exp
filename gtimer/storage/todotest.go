package storage

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/schorlet/exp/gtimer"
)

// TodoTest is a test function.
type TodoTest func(*testing.T, *sqlx.DB, gtimer.TodoStore)

// TodoTester runs a TodoTest function.
type TodoTester func(TodoTest) func(*testing.T)

// TodoTestSuite runs a suite of TodoTest functions.
func TodoTestSuite(t *testing.T, tester TodoTester) {
	t.Run("Todo.Create", tester(todoCreate))
	t.Run("Todo.Read", tester(todoRead))
	t.Run("Todo.Update", tester(todoUpdate))
	t.Run("Todo.Delete", tester(todoDelete))
}

func withID(id string) gtimer.TodoFilter {
	return func(todo *gtimer.Todo) {
		todo.ID = id
	}
}

func withStatus(status string) gtimer.TodoFilter {
	return func(todo *gtimer.Todo) {
		todo.Status = status
	}
}

func todoCreate(t *testing.T, db *sqlx.DB, store gtimer.TodoStore) {
	create1, err := store.Create(db, gtimer.Todo{ID: "st101", Title: "st101"})
	if err != nil {
		t.Fatalf("Unable to create Todo: %v", err)
	}
	if create1.Status != "active" {
		t.Fatalf("Unexpected Todo Status: %s", create1.Status)
	}

	_, err = store.Create(db, gtimer.Todo{ID: "st101", Title: "st101"})
	if err == nil {
		t.Fatal("Expected error when creating Todo")
	}

	create2, err := store.Create(db, gtimer.Todo{Title: "st101"})
	if err != nil {
		t.Fatalf("Unable to create Todo: %v", err)
	}
	if create2.ID == create1.ID {
		t.Fatalf("Unexpected Todo ID: %s", create2.ID)
	}

	todos, err := store.Read(db)
	if err != nil {
		t.Fatalf("Unable to get Todos: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("Unexpected count of Todos: %d", len(todos))
	}
}

func todoRead(t *testing.T, db *sqlx.DB, store gtimer.TodoStore) {
	create, err := store.Create(db, gtimer.Todo{Title: "st101"})
	if err != nil {
		t.Fatalf("Unable to create Todo: %v", err)
	}

	todos, err := store.Read(db, withID(create.ID))
	if err != nil {
		t.Fatalf("Unable to get Todo: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("Unexpected count of Todos: %d", len(todos))
	}

	todos, err = store.Read(db, withID("0"))
	if err == nil {
		t.Fatalf("Unexpected Todo: %v", todos)
	}
	if err != gtimer.ErrNotFound {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("Unexpected count of Todos: %d", len(todos))
	}

	todos, err = store.Read(db, withStatus(create.Status))
	if err != nil {
		t.Fatalf("Unable to get Todo: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("Unexpected count of Todos: %d", len(todos))
	}

	todos, err = store.Read(db, withStatus("foo"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("Unexpected count of Todos: %d", len(todos))
	}
}

func todoUpdate(t *testing.T, db *sqlx.DB, store gtimer.TodoStore) {
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
}

func todoDelete(t *testing.T, db *sqlx.DB, store gtimer.TodoStore) {
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
}
