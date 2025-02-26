package main

import (
	"reflect"
	"strconv"

	"github.com/Bmartin35000/backend-project/todo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func initDatabase() {
	dataSourceName := "host=" + conf.Db.Address + " user=" + conf.Db.User + " password=" + conf.Db.Password +
		" dbname=" + conf.Db.Name + " port=" + strconv.Itoa(conf.Db.Port) + " sslmode=disable"

	var err error
	// Connect to database
	db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // disable default logs
	})
	if err != nil {
		log.WithFields(log.Fields{"dsn": dataSourceName, "details": err.Error()}).Panic("Connection to db failed")
		panic("failed to connect database")
	}
	log.WithFields(log.Fields{"dsn": dataSourceName}).Info("Connection to db successful")

	// Update the db table
	err = db.AutoMigrate(&todo.TodoModel{})
	if err != nil && err.Error() != "simple protocol queries must be run with client_encoding=UTF8" { // ignore one specific error that is not blocking
		log.WithFields(log.Fields{"details": err.Error()}).Panic("Failed to update db schema")
		panic("failed to update db schema")
	}
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
	db.Model(&todoUnique).Where("id = ?", dto.ID).Updates(todo.TodoModel{Title: dto.Title, Completed: dto.Completed})
	if db.Error != nil {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todoUnique), "dto": dto, "details": db.Error}).Error("Update failed")
	} else {
		log.WithFields(log.Fields{"type": reflect.TypeOf(todoUnique), "dto": dto}).Info("Update successful")
	}

	return db.Error
}
