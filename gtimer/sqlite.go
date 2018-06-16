package gtimer

import (
	"github.com/jmoiron/sqlx"
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

// MustConnect connects to the datasource or panics on error.
func MustConnect(datasource string) *DB {
	db := sqlx.MustConnect("sqlite3", datasource)
	return &DB{DB: db}
}

// MustCreateSchema creates the SQL schema or panics on error.
func MustCreateSchema(db *DB) {
	_, err := db.Exec(todoSchema)
	if err != nil {
		panic(err)
	}
}
