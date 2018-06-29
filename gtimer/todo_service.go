package gtimer

// TodoService interface.
type TodoService interface {
	Create(create Todo) (Todo, error)
	Read(filter TodoFilter) (Todos, error)
	Update(update Todo) (Todo, error)
	Delete(id string) error
}

// TodoServiceSQL struct.
type TodoServiceSQL struct {
	DB    *DB
	Todos TodoStore
}

func (service *TodoServiceSQL) Create(create Todo) (Todo, error) {
	return service.Todos.Create(service.DB, create)
}

func (service *TodoServiceSQL) Read(filter TodoFilter) (Todos, error) {
	return service.Todos.Read(service.DB, filter)
}

func (service *TodoServiceSQL) Update(update Todo) (Todo, error) {
	return service.Todos.Update(service.DB, update)
}

func (service *TodoServiceSQL) Delete(id string) error {
	return service.Todos.Delete(service.DB, id)
}
