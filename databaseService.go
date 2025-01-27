package main

import (
	"database/sql"
	"github.com/Bmartin35000/backend-project/todo"
	"os"
	"reflect"
)

var db *sql.DB

func initDatabase() {
	connStr := "postgresql://" + os.Getenv("db.user") + ":" + os.Getenv("db.password") + "@" + os.Getenv("db.address") + ":" +
		os.Getenv("db.port") + "/" + os.Getenv("db.name") + "?sslmode=disable"

	var err error
	// Connect to database
	db, err = sql.Open("postgres", connStr)
	checkError(err)
}

func getDbTodos(sqlQuery string) ([]todo.TodoModel, error) {
	rows, err := db.Query(sqlQuery)
	checkError(err)
	defer func() {
		checkError(rows.Close())
	}()
	var todoModels []todo.TodoModel
	for rows.Next() {
		var todoModel todo.TodoModel
		modelAddresses := reflect.ValueOf(&todoModel).Elem()
		fillModel(modelAddresses, rows)
		todoModels = append(todoModels, todoModel)
	}

	if err != nil {
		return nil, err
	}

	return todoModels, nil
}

func fillModel(modelAddresses reflect.Value, rows *sql.Rows) {
	columnsAddresses := make([]interface{}, modelAddresses.NumField())
	for i := 0; i < modelAddresses.NumField(); i++ {
		columnsAddresses[i] = modelAddresses.Field(i).Addr().Interface()
	}

	checkError(rows.Scan(columnsAddresses...))
}
