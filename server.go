package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bmartin35000/backend-project/todo"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func main() {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
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
		r.Get("/todos", getTodos)
		r.Post("/todo", createTodo)
		r.Put("/todo", updateTodo)
		r.Delete("/todo/{id}", deleteTodo)
	})
	return router
}

func getTodos(response http.ResponseWriter, _ *http.Request) {
	fmt.Println("get todos")

	todoModels, err := getDbTodos()
	if err != nil {
		log.Printf("failed to get todos : %v\n", err.Error())
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not get the todos",
			"error":   err.Error(),
		})
	}
	todoList := todo.MapTodoListModelToDto(todoModels)

	rnd.JSON(response, http.StatusOK, todo.TodoListResponse{
		Message: "All todos retrieved",
		Data:    todoList,
	})
}

func createTodo(response http.ResponseWriter, request *http.Request) {
	fmt.Println("create todo")

	todoDto := todo.TodoDto{}
	if err := json.NewDecoder(request.Body).Decode(&todoDto); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}

	if todoDto.Title == "" {
		log.Println("no title added to response body")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}

	err := createDbTodo(todoDto)
	if err != nil {
		log.Printf("failed to create todo : %v\n", err.Error())
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not create the todo",
			"error":   err.Error(),
		})
	}

	rnd.JSON(response, http.StatusOK, todo.Response{
		Message: "Todo created",
	})
}

func deleteTodo(response http.ResponseWriter, request *http.Request) {
	fmt.Println("delete todo")

	todoId := strings.TrimSpace(chi.URLParam(request, "id"))

	if todoId == "" {
		log.Println("no id added to response body")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "please add an id",
		})
		return
	}

	err := deleteDbTodo(todoId)
	if err != nil {
		log.Printf("failed to delete todo : %v\n", err.Error())
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not delete the todo",
			"error":   err.Error(),
		})
	}

	rnd.JSON(response, http.StatusOK, todo.Response{
		Message: "Todo deleted",
	})
}
func updateTodo(response http.ResponseWriter, request *http.Request) {
	fmt.Println("update todo")

	var todoDto todo.TodoDto
	if err := json.NewDecoder(request.Body).Decode(&todoDto); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}

	if todoDto.Title == "" {
		log.Println("no title added to response body")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}

	err := updateDbTodo(todoDto)
	if err != nil {
		log.Printf("failed to update todo : %v\n", err.Error())
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not update the todo",
			"error":   err.Error(),
		})
	}

	rnd.JSON(response, http.StatusOK, todo.Response{
		Message: "Todo updated",
	})
}
