package mem

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/schorlet/exp/gtimer"
	"github.com/schorlet/exp/gtimer/storage"
)

// TodoStore implements exp/gtimer#TodoStore.
type TodoStore map[string]gtimer.Todo

var _ gtimer.TodoStore = TodoStore{}

func (store TodoStore) Create(_ sqlx.Ext, create gtimer.Todo) (gtimer.Todo, error) {
	if create.ID == "" {
		var err error
		create.ID, err = storage.RandomString(12)
		if err != nil {
			return create, err
		}
	} else if _, err := store.Get(create.ID); err == nil {
		return gtimer.Todo{}, fmt.Errorf("duplicated id: %s", create.ID)
	}
	create.Status = "active"
	create.Created = time.Now()
	create.Updated = time.Now()
	store[create.ID] = create
	return create, nil
}

func (store TodoStore) Read(_ sqlx.Queryer, filter gtimer.TodoFilter) (gtimer.Todos, error) {
	if filter.ID != "" {
		todo, err := store.Get(filter.ID)
		if err != nil {
			return gtimer.Todos{}, err
		}
		return gtimer.Todos{todo}, err
	}
	if filter.Status != "" {
		return store.ByStatus(filter.Status)
	}
	return store.All()
}

func (store TodoStore) Get(id string) (gtimer.Todo, error) {
	if todo, ok := store[id]; ok {
		return todo, nil
	}
	return gtimer.Todo{}, fmt.Errorf("id not found: %s", id)
}

func (store TodoStore) ByStatus(status string) (gtimer.Todos, error) {
	var todos gtimer.Todos
	for _, todo := range store {
		if todo.Status == status {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (store TodoStore) All() (gtimer.Todos, error) {
	todos := make(gtimer.Todos, len(store))
	index := 0
	for _, todo := range store {
		todos[index] = todo
		index++
	}
	return todos, nil
}

func (store TodoStore) Update(_ sqlx.Ext, update gtimer.Todo) (gtimer.Todo, error) {
	todo, err := store.Get(update.ID)
	if err != nil {
		return gtimer.Todo{}, fmt.Errorf("mem update: %v", err)
	}
	if update.Status != "completed" && update.Status != "active" {
		return gtimer.Todo{}, fmt.Errorf("invalid status: %s", update.Status)
	}
	todo.Title = update.Title
	todo.Status = update.Status
	todo.Updated = time.Now()
	store[todo.ID] = todo
	return todo, nil
}

func (store TodoStore) Delete(_ sqlx.Ext, id string) error {
	if _, err := store.Get(id); err != nil {
		return fmt.Errorf("mem delete: %v", err)
	}
	delete(store, id)
	return nil
}
