package todo

import (
	"gorm.io/gorm"
)

type TodoModel struct {
	gorm.Model
	ID        string
	Title     string
	Completed bool
}
