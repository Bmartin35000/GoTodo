package main

import (
	"database/sql"
	"fmt"
	"github.com/Bmartin35000/backend-project/todo"
	"os"
	"reflect"
)

var db *sql.DB

func initDatabase() {
	connStr := "postgresql://" + os.Getenv("db_user") + ":" + os.Getenv("db_password") + "@" + os.Getenv("db_address") + ":" +
		os.Getenv("db_port") + "/" + os.Getenv("db_name") + "?sslmode=disable"
	fmt.Println("Connecting to db : ", connStr)

	var err error
	// Connect to database
	db, err = sql.Open("postgres", connStr)
	checkError(err)
}

func getDbTodos() ([]todo.TodoModel, error) {
	rows, err := db.Query("SELECT * FROM todos")
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

func createDbTodo(dto todo.TodoDto) error {
	sqlQuery := "INSERT INTO todos VALUES (gen_random_uuid (), " + "'" + dto.Title + "'" + ", false, NOW());"
	return executreDbQuery(sqlQuery)
}

func deleteDbTodo(id string) error {
	sqlQuery := "DELETE FROM todos WHERE ID = " + "'" + id + "'" + ";"
	return executreDbQuery(sqlQuery)
}

func updateDbTodo(dto todo.TodoDto) error {
	sqlQuery := "UPDATE todos SET Title = " +
		"'" + dto.Title + "'" +
		", Completed = " +
		toString(dto.Completed) +
		" WHERE ID = " +
		"'" + dto.ID + "'" + ";"
	print(sqlQuery)
	return executreDbQuery(sqlQuery)
}

func executreDbQuery(sqlQuery string) error {
	_, err := db.Exec(sqlQuery)
	checkError(err)

	if err != nil {
		return err
	}

	return nil
}

func toString(bool bool) string {
	if bool {
		return "true"
	} else {
		return "false"
	}
}
