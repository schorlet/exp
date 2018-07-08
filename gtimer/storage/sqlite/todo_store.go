package sqlite

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/schorlet/exp/gtimer"
	"github.com/schorlet/exp/gtimer/storage"
)

const todoSchema = `
	drop index if exists TODO_IDX_STATUS;
	drop table if exists TODO;

	create table TODO (
		ID      text   	  primary key,
		TITLE   text      not null,
		STATUS  text      not null default 'active',
		CREATED datetime  not null default current_timestamp,
		UPDATED datetime  not null default current_timestamp,
		check (STATUS in ('active', 'completed'))
	);

	create index TODO_IDX_STATUS on TODO (STATUS);
`

// TodoStore implements gtimer.TodoStore.
type TodoStore struct {
}

var _ gtimer.TodoStore = TodoStore{}

// MustDefine creates the Todo schema or panics on error.
func (TodoStore) MustDefine(e sqlx.Ext) {
	_, err := e.Exec(todoSchema)
	if err != nil {
		panic(err)
	}
}

// Create handles Todo creation and returns the newly created Todo.
func (store TodoStore) Create(e sqlx.Ext, create gtimer.Todo) (gtimer.Todo, error) {
	query := `
			insert into TODO (ID, TITLE)
			values (?, ?)`

	if create.ID == "" {
		var err error
		create.ID, err = storage.RandomString(12)
		if err != nil {
			return create, err
		}
	}

	_, err := e.Exec(query, create.ID, create.Title)
	if err != nil {
		return create, err
	}

	return store.Get(e, create.ID)
}

// Read returns all Todos with the specified filter.
// Read returns gtimer.ErrNotFound when filtering by ID and when the expected Todo is not found.
// Otherwise the returned Todos may be empty and err be nil.
func (store TodoStore) Read(q sqlx.Queryer, filter gtimer.TodoFilter) (gtimer.Todos, error) {
	if filter.ID != "" {
		todo, err := store.Get(q, filter.ID)
		if err != nil {
			return gtimer.Todos{}, err
		}
		return gtimer.Todos{todo}, err
	}
	if filter.Status != "" {
		return store.ByStatus(q, filter.Status)
	}
	return store.All(q)
}

// Get returns the Todo with the given ID.
func (TodoStore) Get(q sqlx.Queryer, id string) (gtimer.Todo, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where id = ?`

	var todo gtimer.Todo
	err := sqlx.Get(q, &todo, query, id)
	if err == sql.ErrNoRows {
		err = gtimer.ErrNotFound
	}
	return todo, err
}

// ByStatus returns all Todos with the specified Status.
func (TodoStore) ByStatus(q sqlx.Queryer, status string) (gtimer.Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where STATUS = ?
			order by CREATED desc, TITLE asc`

	var todos gtimer.Todos
	err := sqlx.Select(q, &todos, query, status)

	return todos, err
}

// All returns all Todos.
func (TodoStore) All(q sqlx.Queryer) (gtimer.Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			order by CREATED desc, TITLE asc`

	var todos gtimer.Todos
	err := sqlx.Select(q, &todos, query)

	return todos, err
}

// Update updates the Title and Status of the Todo with the given ID.
func (store TodoStore) Update(e sqlx.Ext, update gtimer.Todo) (gtimer.Todo, error) {
	query := `
			update TODO set TITLE = ?,
							STATUS = ?,
							UPDATED = current_timestamp
			where ID = ?`

	r, err := e.Exec(query, update.Title, update.Status, update.ID)
	if err != nil {
		return update, err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		return update, gtimer.ErrNotFound
	}

	return store.Get(e, update.ID)
}

// Delete deletes the Todo with the given ID.
func (TodoStore) Delete(e sqlx.Ext, id string) error {
	query := `delete from TODO where ID = ?`

	r, err := e.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		err = gtimer.ErrNotFound
	}
	return err
}
