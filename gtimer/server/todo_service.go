package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/schorlet/exp/gtimer"
)

// TodoService implements gtimer.TodoService.
type TodoService struct {
	DB    *sqlx.DB
	Store gtimer.TodoStore
}

var _ gtimer.TodoService = new(TodoService)

// Create handles Todo creation and returns the newly created Todo.
func (todos *TodoService) Create(create gtimer.Todo) (gtimer.Todo, error) {
	return todos.Store.Create(todos.DB, create)
}

// Read searches for Todos according to the specified filter.
func (todos *TodoService) Read(filter gtimer.TodoFilter) (gtimer.Todos, error) {
	return todos.Store.Read(todos.DB, filter)
}

// Update handles Todo modification and returns the updated Todo.
func (todos *TodoService) Update(update gtimer.Todo) (gtimer.Todo, error) {
	return todos.Store.Update(todos.DB, update)
}

// Delete handles Todo deletion.
func (todos *TodoService) Delete(id string) error {
	return todos.Store.Delete(todos.DB, id)
}
