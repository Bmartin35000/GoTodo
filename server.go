package main

import (
	"fmt"
	"github.com/Bmartin35000/backend-project/todo"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Mount("/", todoHandlers())

	server := &http.Server{
		Addr:         ":" + os.Getenv("server.port"),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// start the server
	fmt.Println("Server started on port", os.Getenv("server.port"))
	if err := server.ListenAndServe(); err != nil {
		log.Printf("listen:%s\n", err)

	}
}

func init() {
	rnd = renderer.New()
}

func todoHandlers() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		//r.Post("/", createTodo)
		//r.Put("/{id}", updateTodo)
		//r.Delete("/{id}", deleteTodo)
	})
	return router
}

func getTodos(rw http.ResponseWriter, _ *http.Request) {
	fmt.Println("get todos")

	sqlQuery := "SELECT * FROM todos"
	todoModels, err := getDbTodos(sqlQuery)
	if err != nil {
		log.Printf("failed to fetch todo records from the db: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Could not fetch the todo collection",
			"error":   err.Error(),
		})
	}
	todoList := todo.MapTodoListModelToDto(todoModels)

	rnd.JSON(rw, http.StatusOK, todo.GetTodoResponse{
		Message: "All todos retrieved",
		Data:    todoList,
	})
}
