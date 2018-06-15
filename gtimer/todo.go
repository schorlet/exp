package gtimer

import (
	"fmt"
	"time"
)

// Todo struct
type Todo struct {
	ID      string    `db:"ID"`
	Title   string    `db:"TITLE" `
	Status  string    `db:"STATUS"`
	Created time.Time `db:"CREATED"`
	Updated time.Time `db:"UPDATED"`
}

// Todos slice
type Todos []Todo

func (t Todo) String() string {
	return fmt.Sprintf("Todo{ID:%s, Title:%s, Status:%s, Created:%s, Updated:%s}",
		t.ID, t.Title, t.Status, t.Created, t.Updated)
}
