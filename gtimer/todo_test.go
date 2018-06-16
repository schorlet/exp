package gtimer

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func withDB(fn func(db *DB)) {
	db := MustConnect(":memory:")
	defer db.Close()
	MustCreateSchema(db)

	fn(db)
}

func TestCreateTodo(t *testing.T) {
	withDB(func(db *DB) {
		todo1, err := CreateTodo(db, "st101")
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		todos, err := GetTodos(db)
		if err != nil {
			t.Fatalf("Unable to get Todos: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todo2, err := CreateTodo(db, "st101")
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}
		if todo2.ID == todo1.ID {
			t.Fatalf("Unexpected Todo ID: %s", todo2.ID)
		}

		todos, err = GetTodos(db)
		if err != nil {
			t.Fatalf("Unable to get Todos: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestGetTodo(t *testing.T) {
	withDB(func(db *DB) {
		todo, err := CreateTodo(db, "st101")
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		todo, err = GetTodo(db, todo.ID)
		if err != nil {
			t.Fatalf("Unable to get Todo: %v", err)
		}

		todo, err = GetTodo(db, "1")
		if err == nil {
			t.Fatalf("Unexpected Todo: %v", todo)
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	withDB(func(db *DB) {
		todo, err := CreateTodo(db, "st101")
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		err = DeleteTodo(db, todo.ID)
		if err != nil {
			t.Fatalf("Unable to delete Todo: %v", err)
		}

		err = DeleteTodo(db, "1")
		if err == nil {
			t.Fatalf("Expected error when deleting Todo")
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	withDB(func(db *DB) {
		todo, err := CreateTodo(db, "st101")
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		todos, err := GetTodosByStatus(db, todo.Status)
		if err != nil {
			t.Fatalf("Unable to get Todos by Status: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todo.Status = "completed"
		err = UpdateTodo(db, todo)
		if err != nil {
			t.Fatalf("Unable to update Todo: %v", err)
		}

		todos, err = GetTodosByStatus(db, todo.Status)
		if err != nil {
			t.Fatalf("Unable to get Todos by Status: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todo.Status = "foo"
		err = UpdateTodo(db, todo)
		if err == nil {
			t.Fatal("Expected error when updating Todo")
		}

		todos, err = GetTodosByStatus(db, todo.Status)
		if err != nil {
			t.Fatalf("Unable to get Todos by Status: %v", err)
		}
		if len(todos) != 0 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}
