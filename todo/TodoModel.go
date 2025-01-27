package todo

import "time"

type TodoModel struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
}
