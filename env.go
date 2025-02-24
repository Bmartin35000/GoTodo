package main

import (
	"os"
)

func init() {
	setEnvVariableIfNotExist("server.port", "80")
	setEnvVariableIfNotExist("db_user", "postgres")
	setEnvVariableIfNotExist("db_password", "postgres")
	setEnvVariableIfNotExist("db_address", "localhost")
	setEnvVariableIfNotExist("db_port", "55555")
	setEnvVariableIfNotExist("db_name", "postgres")
	setEnvVariableIfNotExist("environment", "development")
	initDatabase() // force init orders due to strong connection
}

func setEnvVariableIfNotExist(key string, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}
