package main

import (
	"os"
)

func init() {
	os.Setenv("server.port", "8080")
	os.Setenv("db.user", "postgres")
	os.Setenv("db.password", "postgres")
	os.Setenv("db.address", "localhost")
	os.Setenv("db.port", "55555")
	os.Setenv("db.name", "todos")
	initDatabase() // force init orders due to strong connection
}
