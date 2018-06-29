package gtimer

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func withService(fn func(TodoService)) {
	db := MustConnect(":memory:")
	defer db.Close()
	MustCreateSchema(db)

	service := TodoServiceSQL{DB: db, Todos: TodoSqlite{}}
	fn(&service)
}

func TestCreateTodo(t *testing.T) {
	withService(func(service TodoService) {
		create1, err := service.Create(Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		create2, err := service.Create(Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}
		if create2.ID == create1.ID {
			t.Fatalf("Unexpected Todo ID: %s", create2.ID)
		}

		todos, err := service.Read(TodoFilter{})
		if err != nil {
			t.Fatalf("Unable to get Todos: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestGetTodo(t *testing.T) {
	withService(func(service TodoService) {
		create, err := service.Create(Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		todos, err := service.Read(TodoFilter{ID: create.ID})
		if err != nil {
			t.Fatalf("Unable to get Todo: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}

		todos, err = service.Read(TodoFilter{ID: "0"})
		if err == nil {
			t.Fatalf("Unexpected Todo: %v", todos)
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	withService(func(service TodoService) {
		create, err := service.Create(Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		err = service.Delete(create.ID)
		if err != nil {
			t.Fatalf("Unable to delete Todo: %v", err)
		}

		err = service.Delete("0")
		if err == nil {
			t.Fatalf("Expected error when deleting Todo")
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	withService(func(service TodoService) {
		create, err := service.Create(Todo{Title: "st101"})
		if err != nil {
			t.Fatalf("Unable to create Todo: %v", err)
		}

		create.Status = "completed"
		update, err := service.Update(create)
		if err != nil {
			t.Fatalf("Unable to update Todo: %v", err)
		}

		update.Status = "foo"
		update, err = service.Update(update)
		if err == nil {
			t.Fatal("Expected error when updating Todo")
		}
	})
}
