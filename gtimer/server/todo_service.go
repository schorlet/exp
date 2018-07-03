package server

import (
	"github.com/schorlet/exp/gtimer"
	"github.com/schorlet/exp/sql"
)

// TodoService implements #gtimer.TodoService.
type TodoService struct {
	DB    *sql.DB
	Store gtimer.TodoStore
}

var _ gtimer.TodoService = new(TodoService)

func (todos *TodoService) Create(create gtimer.Todo) (gtimer.Todo, error) {
	return todos.Store.Create(todos.DB, create)
}

func (todos *TodoService) Read(filter gtimer.TodoFilter) (gtimer.Todos, error) {
	return todos.Store.Read(todos.DB, filter)
}

func (todos *TodoService) Update(update gtimer.Todo) (gtimer.Todo, error) {
	return todos.Store.Update(todos.DB, update)
}

func (todos *TodoService) Delete(id string) error {
	return todos.Store.Delete(todos.DB, id)
}
