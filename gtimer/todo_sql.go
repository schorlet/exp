package gtimer

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"

	"github.com/jmoiron/sqlx"
)

// TodoStore interface.
type TodoStore interface {
	Create(e sqlx.Ext, create Todo) (Todo, error)
	Read(q sqlx.Queryer, filter TodoFilter) (Todos, error)
	Update(e sqlx.Ext, update Todo) (Todo, error)
	Delete(e sqlx.Ext, id string) error
}

// TodoSqlite type.
type TodoSqlite struct {
}

// Create creates a Todo.
func (sqlite TodoSqlite) Create(e sqlx.Ext, create Todo) (Todo, error) {
	query := `
			insert into TODO (ID, TITLE)
			values (?, ?)`

	if create.ID == "" {
		var err error
		create.ID, err = randomString(12)
		if err != nil {
			return create, err
		}
	}

	_, err := e.Exec(query, create.ID, create.Title)
	if err != nil {
		return create, err
	}

	return sqlite.Get(e, create.ID)
}

func randomString(length int) (string, error) {
	// https://www.commandlinefu.com/commands/view/24071/generate-random-text-based-on-length
	raw := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, raw)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(raw)[:length], nil
}

// Read returns all Todos with the specified filter.
func (sqlite TodoSqlite) Read(q sqlx.Queryer, filter TodoFilter) (Todos, error) {
	if filter.ID != "" {
		todo, err := sqlite.Get(q, filter.ID)
		return Todos{todo}, err
	}
	if filter.Status != "" {
		return sqlite.ByStatus(q, filter.Status)
	}
	return sqlite.All(q)
}

// Get returns the Todo with the given ID.
func (TodoSqlite) Get(q sqlx.Queryer, id string) (Todo, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where id = ?`

	var todo Todo
	err := sqlx.Get(q, &todo, query, id)
	return todo, err
}

// ByStatus returns all Todos with the specified Status.
func (TodoSqlite) ByStatus(q sqlx.Queryer, status string) (Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where STATUS = ?
			order by CREATED desc, TITLE asc`

	var todos Todos
	err := sqlx.Select(q, &todos, query, status)

	return todos, err
}

// All returns all Todos.
func (TodoSqlite) All(q sqlx.Queryer) (Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			order by CREATED desc, TITLE asc`

	var todos Todos
	err := sqlx.Select(q, &todos, query)

	return todos, err
}

// Update updates the Title and Status of the Todo with the given ID.
func (sqlite TodoSqlite) Update(e sqlx.Ext, update Todo) (Todo, error) {
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
		return update, sql.ErrNoRows
	}

	return sqlite.Get(e, update.ID)
}

// Delete deletes the Todo with the given ID.
func (TodoSqlite) Delete(e sqlx.Ext, id string) error {
	query := `delete from TODO where ID = ?`

	r, err := e.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		err = sql.ErrNoRows
	}
	return err
}
