package main

import (
	"os"

	"reflect"

	"github.com/Bmartin35000/backend-project/todo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func initDatabase() {
	dataSourceName := "host=" + os.Getenv("db_address") + " user=" + os.Getenv("db_user") + " password=" + os.Getenv("db_password") +
		" dbname=" + os.Getenv("db_name") + " port=" + os.Getenv("db_port") + " sslmode=disable"

	var err error
	// Connect to database
	db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // disable default logs
	})
	if err != nil {
		log.WithFields(log.Fields{"dsn": dataSourceName}).Panic("Connection to db failed")
		panic("failed to connect database")
	}
	log.WithFields(log.Fields{"dsn": dataSourceName}).Info("Connection to db successful")

	// Update the db table
	db.AutoMigrate(&todo.TodoModel{})
}

func getDbTodos() ([]todo.TodoModel, error) {
	var todos []todo.TodoModel
	db.Find(&todos)
	if db.Error != nil {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todos), "details": db.Error}).Error("Get failed")
	} else {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todos), "length": len(todos)}).Info("Get successful")
	}

	return todos, db.Error
}

func createDbTodo(dto todo.TodoDto) error {
	id := uuid.New().String()
	db.Create(&todo.TodoModel{ID: id, Title: dto.Title})
	if db.Error != nil {
		log.WithFields(log.Fields{"type": reflect.TypeOf(dto), "dto": dto, "details": db.Error}).Error("Create failed")
	} else {
		log.WithFields(log.Fields{"type": reflect.TypeOf(dto), "dto": dto}).Info("Create successful")
	}

	return db.Error
}

func deleteDbTodo(id string) error {
	var todoUnique todo.TodoModel
	db.Delete(&todoUnique, "id = ?", id)
	if db.Error != nil {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todoUnique), "id": id, "details": db.Error}).Error("Delete failed")
	} else {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todoUnique), "id": id}).Info("Delete successful")
	}

	return db.Error
}

func updateDbTodo(dto todo.TodoDto) error {
	var todoUnique todo.TodoModel
	db.Model(&todoUnique).Updates(todo.TodoModel{Title: dto.Title, Completed: dto.Completed})
	if db.Error != nil {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todoUnique), "dto": dto, "details": db.Error}).Error("Update failed")
	} else {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todoUnique), "dto": dto}).Info("Update successful")
	}

	return db.Error
}
