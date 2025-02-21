# Prerequisite
create a postgres db, change credentials in env.go

# launch project
go run .

# Sources
source project : https://www.agirlcodes.dev/build-todo-app-backend-golang-tutorial

postgres connection : https://blog.logrocket.com/building-simple-app-go-postgresql/

# Docker
## Create the image
docker build -t bamartin35/go-server-todo:1.0 .

## Run the back-end container
docker run -p 8080:80 bamartin35/go-server-todo:1.0

## Run the whole application
docker compose up

## Monitoring database
host: localhost
port:55555
dbname: postgres
username: postgres
password: postgres

## Accessing the APIs
address: host.docker.internal:8080

