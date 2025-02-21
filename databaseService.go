package main

import (
	"fmt"
	"github.com/Bmartin35000/backend-project/todo"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func initDatabase() {
	dataSourceName := "host=" + os.Getenv("db_address") + " user=" + os.Getenv("db_user") + " password=" + os.Getenv("db_password") +
		" dbname=" + os.Getenv("db_name") + " port=" + os.Getenv("db_port") + " sslmode=disable"
	fmt.Println("Connecting to db : ", dataSourceName)

	var err error
	// Connect to database
	db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Update the db table
	db.AutoMigrate(&todo.TodoModel{})
}

func getDbTodos() ([]todo.TodoModel, error) {
	var todos []todo.TodoModel
	db.Find(&todos)
	checkError(db.Error)
	fmt.Println("len todos : ", len(todos))

	return todos, db.Error
}

func createDbTodo(dto todo.TodoDto) error {
	id := uuid.New().String()
	db.Create(&todo.TodoModel{ID: id, Title: dto.Title})
	checkError(db.Error)

	return db.Error
}

func deleteDbTodo(id string) error {
	var todoUnique []todo.TodoModel
	db.First(&todoUnique, "id = ?", id)
	db.Delete(&todoUnique, "id = ?", id)
	checkError(db.Error)

	return db.Error
}

func updateDbTodo(dto todo.TodoDto) error {
	var todoUnique []todo.TodoModel
	db.First(&todoUnique, "id = ?", dto.ID)
	db.Model(&todoUnique).Updates(todo.TodoModel{Title: dto.Title, Completed: dto.Completed})
	checkError(db.Error)

	return db.Error
}
